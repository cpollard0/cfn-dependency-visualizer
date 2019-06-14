package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"fmt"
	"os"
	"strings"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"net/http"
)


type cloudFormationNode struct
{
	Id int `json:"id"`
	StackName string  `json:"name"`
}

type cloudFormationLink struct
{
	Source string `json:"source"`
	Target string `json:"target"`
}

type jsonOutput struct
{
	Nodes [] cloudFormationNode `json:"nodes"`
	Links [] cloudFormationLink `json:"links"`
}
// Initialize a session that the SDK uses to load
// credentials from the shared credentials file ~/.aws/credentials
// and configuration from the shared configuration file ~/.aws/config.
var sess = session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))

    // Create CloudFormation client in region
var	svc = cloudformation.New(sess, &aws.Config{
		Region: aws.String("us-east-1"),
	})

var counter = 1
var cfNodes []cloudFormationNode

func main() {
	var finalOutput jsonOutput
	fmt.Println("begin processing")
	var cfLinks []cloudFormationLink
	var listOfStrings []cloudformation.Export

	input := &cloudformation.ListExportsInput{}

	for {
		listExportsOutputs, err := svc.ListExports(input)
		if err != nil {
			fmt.Println("Got error listing exports:")
			fmt.Println(err.Error())
			os.Exit(1)
		}
		for _, stack := range listExportsOutputs.Exports {
			listOfStrings = append(listOfStrings,*stack)
		}
		if listExportsOutputs.NextToken == nil {
			break
		} else {
			input.NextToken=listExportsOutputs.NextToken
		}
	}
	fmt.Println("after exports")

    for _, stack := range listOfStrings {
		s := strings.Split(*stack.ExportingStackId, "/")

		listImportsInput := &cloudformation.ListImportsInput{ExportName: stack.Name}
		listImportsResponse, err := svc.ListImports(listImportsInput)
		cfNode := cloudFormationNode{counter,s[1]}
		cfNodes=AppendNodeIfMissing(cfNodes, cfNode)
		if err == nil {
			for _, imports := range listImportsResponse.Imports{
				cfLink := cloudFormationLink{s[1],*imports}
				// even if something is not exporting anything but is inputing something
				// it need sto be a node
				linkFromNode := cloudFormationNode{counter,cfLink.Target}
				cfNodes=AppendNodeIfMissing(cfNodes, linkFromNode)
				cfLinks=AppendIfMissing(cfLinks, cfLink)
			}
		}
	}

	// set the ID on all the nodes
	for i := range cfNodes {
		cfNodes[i].Id = counter
		counter++
	}
	for i := range cfLinks {
		cfLinks[i].Source = lookupSourceByName(cfLinks[i].Source)
		cfLinks[i].Target = lookupSourceByName(cfLinks[i].Target)
		counter++
	}
	fmt.Println("end processing")
	finalOutput.Nodes = cfNodes
	finalOutput.Links = cfLinks
	marshalledJson, _ := json.Marshal(finalOutput)

	ioutil.WriteFile("output.json", marshalledJson, 0644)
	// start webserver
	http.Handle("/", http.FileServer(http.Dir("./")))
	http.ListenAndServe(":8080", nil)

}
func AppendIfMissing(slice []cloudFormationLink, i cloudFormationLink) []cloudFormationLink {
    for _, ele := range slice {
        if ele == i {
            return slice
		}
	}
    return append(slice, i)
}
func lookupSourceByName(name string) string{
	for _, export := range cfNodes {
		if name == export.StackName {
			return strconv.Itoa(export.Id)
		}
	}
	return "0"
}
func AppendNodeIfMissing(slice []cloudFormationNode, i cloudFormationNode) []cloudFormationNode {
    for _, ele := range slice {
        if ele.StackName == i.StackName {
            return slice
        }
    }
    return append(slice, i)
}
