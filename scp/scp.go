package scp

import (
	"fmt"
	"io"
	"log"
	"os"
	// "path/filepath" // Removed: "path/filepath" imported and not used
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

// SCPClient represents an SCP client connection for file transfers.
type SCPClient struct {
	Host     string
	User     string
	Password string
	sshClient *ssh.Client
	mu       sync.Mutex // Mutex to protect sshClient access
}

// NewSCPClient creates a new SCP client.
func NewSCPClient(host, user, password string) (*SCPClient, error) {
	return &SCPClient{
		Host:     host,
		User:     user,
		Password: password,
	}, nil
}

// Connect establishes an SSH connection if one is not already open.
func (c *SCPClient) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// If client is already connected and active, return
	if c.sshClient != nil && c.sshClient.Conn != nil && c.sshClient.Conn.LocalAddr() != nil {
		return nil
	}

	hostKeyCallback := ssh.InsecureIgnoreHostKey() // WARNING: Insecure for production

	config := &ssh.ClientConfig{
		User: c.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(c.Password),
		},
		HostKeyCallback: hostKeyCallback,
		Timeout: 15 * time.Second,
	}

	client, err := ssh.Dial("tcp", c.Host+":22", config)
	if err != nil {
		return fmt.Errorf("failed to dial SSH: %w", err)
	}
	c.sshClient = client
	log.Printf("SSH client connected to %s", c.Host)
	return nil
}

// Close closes the SSH connection.
func (c *SCPClient) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.sshClient != nil {
		err := c.sshClient.Close()
		if err != nil {
			log.Printf("Error closing SSH client: %v", err)
		}
		c.sshClient = nil // Clear client after closing
		log.Printf("SSH client connection to %s closed", c.Host)
	}
}

// CheckRemoteDirectory attempts to determine if the target directory exists on the remote.
func (c *SCPClient) CheckRemoteDirectory(targetFileDir string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.sshClient == nil || c.sshClient.Conn == nil || c.sshClient.Conn.LocalAddr() == nil {
		return "", fmt.Errorf("SSH client not connected. Call Connect() first.")
	}

	session, err := c.sshClient.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create SSH session for directory check: %w", err)
	}
	defer session.Close()

	cmd := fmt.Sprintf("test -d %s", strings.ReplaceAll(targetFileDir, "\\", "/"))
	if err := session.Run(cmd); err != nil {
		return "", fmt.Errorf("remote directory check failed: %s not found or is not a directory. Error: %w", targetFileDir, err)
	}

	log.Printf("Found remote directory: %s", targetFileDir)
	return targetFileDir, nil
}

// UploadFile uploads content from an io.Reader to a remote path using raw SCP commands over SSH.
// It requires the fileName, fileSize, and fileMode for the SCP protocol header.
func (c *SCPClient) UploadFile(reader io.Reader, remotePath string, fileName string, fileSize int64, fileMode os.FileMode) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	log.Printf("Attempting to upload %s (size: %d) to %s:%s using SCP", fileName, fileSize, c.Host, remotePath)

	if c.sshClient == nil || c.sshClient.Conn == nil || c.sshClient.Conn.LocalAddr() == nil {
		return fmt.Errorf("SSH client not connected. Call Connect() first.")
	}

	session, err := c.sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session for SCP upload: %w", err)
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe for SCP: %w", err)
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	remoteScpCommand := fmt.Sprintf("scp -t %s", strings.ReplaceAll(remotePath, "\\", "/"))

	if err := session.Start(remoteScpCommand); err != nil {
		return fmt.Errorf("failed to start remote SCP command: %w", err)
	}

	// Send file header: "C<mode> <length> <filename>\n"
	fmt.Fprintf(stdin, "C%#o %d %s\n", fileMode, fileSize, fileName)

	bytesWritten, err := io.Copy(stdin, reader)
	if err != nil {
		return fmt.Errorf("failed to write file content to SCP: %w", err)
	}

	if bytesWritten != fileSize {
		log.Printf("Warning: Bytes written (%d) does not match expected file size (%d) for %s", bytesWritten, fileSize, fileName)
	}

	fmt.Fprint(stdin, "\x00") // Send null byte to indicate end of file content

	stdin.Close() // Close stdin to signal end of data

	if err := session.Wait(); err != nil {
		return fmt.Errorf("remote SCP command failed: %w", err)
	}

	log.Printf("Successfully uploaded %s (%d bytes) to %s using SCP", fileName, bytesWritten, remotePath)
	return nil
}


