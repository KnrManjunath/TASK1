package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// Get the kubeconfig file path
	kubeconfig := flag.String("kubeconfig", filepath.Join(os.Getenv("HOME"), ".kube", "config"), "absolute path to the kubeconfig file")
	flag.Parse()

	// Build the clientset
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Fetch and print pods
	pods, err := getPods(clientset)
	if err != nil {
		panic(err.Error())
	}
	printPods(pods)

	time.Sleep(5 * time.Second) // Add a delay of 5 seconds

	// Get pod details
	if len(pods) > 0 {
		getPodDetails(clientset, pods[0])
	}

	time.Sleep(2 * time.Second) // Add a delay of 2 seconds

	// Delete a pod by ID
	if len(pods) > 0 {
		podID := pods[0].Name
		deletePodByID(clientset, podID)
	}

	// Fetch and print pods after deletion
	updatedPods, err := getPods(clientset)
	if err != nil {
		panic(err.Error())
	}
	printPods(updatedPods)
}

// Function to fetch pods
func getPods(clientset *kubernetes.Clientset) ([]corev1.Pod, error) {
	fmt.Println("Fetching pods...")
	pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return pods.Items, nil
}

// Function to print pods
func printPods(pods []corev1.Pod) {
	fmt.Println("Pods:")
	for _, pod := range pods {
		fmt.Printf("Name: %s, Status: %s\n", pod.Name, pod.Status.Phase)
	}
}

// Function to get pod details
func getPodDetails(clientset *kubernetes.Clientset, pod corev1.Pod) {
	fmt.Println("Fetching pod details...")
	fetchedPod, err := clientset.CoreV1().Pods("default").Get(context.TODO(), pod.Name, metav1.GetOptions{})
	if err != nil {
		// Handle the error appropriately for your use case
		fmt.Printf("Error fetching pod details: %v\n", err)
		return
	}
	fmt.Printf("Pod details:\nName: %s\nStatus: %s\n", fetchedPod.Name, fetchedPod.Status.Phase)
}

// Function to delete a pod by ID
func deletePodByID(clientset *kubernetes.Clientset, podID string) {
	fmt.Println("Deleting pod...")
	err := clientset.CoreV1().Pods("default").Delete(context.TODO(), podID, metav1.DeleteOptions{})
	if err != nil {
		// Handle the error appropriately for your use case
		fmt.Printf("Error deleting pod: %v\n", err)
		return
	}
	fmt.Printf("Pod %s deleted successfully.\n", podID)
}

// Fetching pods...
// Pods:
// Name: pod1, Status: Running
// Name: pod2, Status: Pending
// Fetching pod details...
// Pod details:
// Name: pod1
// Status: Running
// Deleting pod...
// Pod pod1 deleted successfully.
// Fetching pods...
// Pods:
// Name: pod2, Status: Pending
