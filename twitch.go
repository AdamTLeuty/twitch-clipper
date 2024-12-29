package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// Define a struct to map the JSON response.
type Metadata struct {
	Author string `json:"author"`
}

type StreamlinkResponse struct {
	Metadata Metadata `json:"metadata"`
}

// StreamlinkChannelName fetches the channel name from a given Twitch URL using streamlink
func StreamlinkChannelName(twitchURL string) (string, error) {
	// Execute the streamlink command with the --json flag
	cmd := exec.Command("streamlink", "--json", twitchURL)

	// Get the command output
	output, err := cmd.Output()
	if err != nil {
		return "", nil
	}

	/// Parse the JSON output.
	var response StreamlinkResponse
	if err := json.Unmarshal(output, &response); err != nil {
		return "", nil
	}

	return response.Metadata.Author, nil
}

// RecordTwitchStream takes a Twitch URL and records a 10-second clip.
func RecordTwitchStream(url string) error {
	// Sanitize the user-submitted URL
	sanitizedURL, err := SanitizeURL(url)
	if err != nil {
		return fmt.Errorf("failed to sanitize URL: %w", err)
	}

	// Define the Streamlink command
	streamlinkArgs := []string{
		"--twitch-disable-ads", // Disable ads
		sanitizedURL,           // The Twitch URL
		"best",                 // Quality (best stream)
		"-o",                   // Output flag (to pipe the stream)
		"-",                    // Output to stdout
	}

	channelName, err := StreamlinkChannelName(sanitizedURL)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Channel Name: %s\n", channelName)
	}

	output := "inClips/" + channelName + ".mp4"

	// Define the FFmpeg command
	ffmpegArgs := []string{
		"-i", "pipe:", // Read from stdin (pipe)
		"-t", "10", // Duration of the clip (10 seconds)
		"-c", "copy", // Use stream copy mode for efficiency
		output, // Output file name
	}

	// Create a pipe to connect the Streamlink output to FFmpeg input
	pipeReader, pipeWriter := io.Pipe()

	// Set up the Streamlink command
	streamlinkCmd := exec.Command("streamlink", streamlinkArgs...)
	streamlinkCmd.Stdout = pipeWriter
	streamlinkCmd.Stderr = os.Stderr

	// Set up the FFmpeg command
	ffmpegCmd := exec.Command("ffmpeg", ffmpegArgs...)
	ffmpegCmd.Stdin = pipeReader
	ffmpegCmd.Stderr = os.Stderr

	// Start the Streamlink command
	err = streamlinkCmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start Streamlink: %w", err)
	}

	// Start the FFmpeg command
	err = ffmpegCmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start FFmpeg: %w", err)
	}

	// Wait for Streamlink to finish
	err = streamlinkCmd.Wait()
	if err != nil {
		return fmt.Errorf("Streamlink command failed: %w", err)
	}

	// Wait for FFmpeg to finish
	err = ffmpegCmd.Wait()
	if err != nil {
		return fmt.Errorf("FFmpeg command failed: %w", err)
	}

	// Close the pipe once done
	pipeWriter.Close()
	pipeReader.Close()

	return nil
}
