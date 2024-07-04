package p

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/compute/metadata"
	"google.golang.org/api/compute/v1"
)

func startInstance(w http.ResponseWriter, r *http.Request) {
	// Project and Zone configuration
	projectID := os.Getenv("project_id")
	zone := os.Getenv("zone")
	instanceName := os.Getenv("vm_instance")

	// Start the VM
	ctx := context.Background()
	computeService, err := compute.NewService(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating compute service: %v", err), http.StatusInternalServerError)
		return
	}

	_, err = computeService.Instances.Start(projectID, zone, instanceName).Context(ctx).Do()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error starting instance: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Println("Instance started successfully")

	// Wait for the server to be ready
	for {
		if instanceIsReady(projectID, zone, instanceName) {
			break
		}
		fmt.Println("Server not ready, waiting 1 second...")
		time.Sleep(1 * time.Second)
	}

	// Get the server's external IP
	serverIP, err := getServerIP(instanceName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting server IP: %v", err), http.StatusInternalServerError)
		return
	}

	// Get the caller's IP
	callerIP := r.Header.Get("X-Forwarded-For")
	if callerIP == "" {
		callerIP = r.RemoteAddr
	}

	// Create a firewall rule
	firewallRuleName := fmt.Sprintf("minecraft-fw-rule-%d", time.Now().Unix())
	firewallRule := &compute.Firewall{
		Name: firewallRuleName,
		Allowed: []*compute.FirewallAllowed{
			{
				IPProtocol: "tcp",
				Ports:      []string{"25565"},
			},
		},
		SourceRanges: []string{fmt.Sprintf("%s/32", callerIP)},
		TargetTags:   []string{"minecraft-server"},
	}

	_, err = computeService.Firewalls.Insert(projectID, firewallRule).Context(ctx).Do()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating firewall rule: %v", err), http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf("Minecraft Server Started! <br />The IP address of the Minecraft server is: %s:25565<br />Your IP address is %s<br />A Firewall rule named %s has been created for you.", serverIP, callerIP, firewallRuleName)
	w.Write([]byte(response))
}

func getServerIP(instanceName string) (string, error) {
	// Using the metadata server to get the external IP
	// This is a more reliable way in the context of Cloud Functions
	return metadata.Get("instance/network-interfaces/0/access-configs/0/external-ip")
}

func instanceIsReady(projectID, zone, instanceName string) bool {
	// This would typically involve checking the instance status
	// or perhaps making a connection attempt to the Minecraft server
	return true // Placeholder for a more robust check
}
