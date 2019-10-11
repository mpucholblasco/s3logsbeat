// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package ec2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awsutil"
)

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DescribeSecurityGroupsRequest
type DescribeSecurityGroupsInput struct {
	_ struct{} `type:"structure"`

	// Checks whether you have the required permissions for the action, without
	// actually making the request, and provides an error response. If you have
	// the required permissions, the error response is DryRunOperation. Otherwise,
	// it is UnauthorizedOperation.
	DryRun *bool `locationName:"dryRun" type:"boolean"`

	// The filters. If using multiple filters for rules, the results include security
	// groups for which any combination of rules - not necessarily a single rule
	// - match all filters.
	//
	//    * description - The description of the security group.
	//
	//    * egress.ip-permission.cidr - An IPv4 CIDR block for an outbound security
	//    group rule.
	//
	//    * egress.ip-permission.from-port - For an outbound rule, the start of
	//    port range for the TCP and UDP protocols, or an ICMP type number.
	//
	//    * egress.ip-permission.group-id - The ID of a security group that has
	//    been referenced in an outbound security group rule.
	//
	//    * egress.ip-permission.group-name - The name of a security group that
	//    has been referenced in an outbound security group rule.
	//
	//    * egress.ip-permission.ipv6-cidr - An IPv6 CIDR block for an outbound
	//    security group rule.
	//
	//    * egress.ip-permission.prefix-list-id - The ID (prefix) of the AWS service
	//    to which a security group rule allows outbound access.
	//
	//    * egress.ip-permission.protocol - The IP protocol for an outbound security
	//    group rule (tcp | udp | icmp or a protocol number).
	//
	//    * egress.ip-permission.to-port - For an outbound rule, the end of port
	//    range for the TCP and UDP protocols, or an ICMP code.
	//
	//    * egress.ip-permission.user-id - The ID of an AWS account that has been
	//    referenced in an outbound security group rule.
	//
	//    * group-id - The ID of the security group.
	//
	//    * group-name - The name of the security group.
	//
	//    * ip-permission.cidr - An IPv4 CIDR block for an inbound security group
	//    rule.
	//
	//    * ip-permission.from-port - For an inbound rule, the start of port range
	//    for the TCP and UDP protocols, or an ICMP type number.
	//
	//    * ip-permission.group-id - The ID of a security group that has been referenced
	//    in an inbound security group rule.
	//
	//    * ip-permission.group-name - The name of a security group that has been
	//    referenced in an inbound security group rule.
	//
	//    * ip-permission.ipv6-cidr - An IPv6 CIDR block for an inbound security
	//    group rule.
	//
	//    * ip-permission.prefix-list-id - The ID (prefix) of the AWS service from
	//    which a security group rule allows inbound access.
	//
	//    * ip-permission.protocol - The IP protocol for an inbound security group
	//    rule (tcp | udp | icmp or a protocol number).
	//
	//    * ip-permission.to-port - For an inbound rule, the end of port range for
	//    the TCP and UDP protocols, or an ICMP code.
	//
	//    * ip-permission.user-id - The ID of an AWS account that has been referenced
	//    in an inbound security group rule.
	//
	//    * owner-id - The AWS account ID of the owner of the security group.
	//
	//    * tag:<key> - The key/value combination of a tag assigned to the resource.
	//    Use the tag key in the filter name and the tag value as the filter value.
	//    For example, to find all resources that have a tag with the key Owner
	//    and the value TeamA, specify tag:Owner for the filter name and TeamA for
	//    the filter value.
	//
	//    * tag-key - The key of a tag assigned to the resource. Use this filter
	//    to find all resources assigned a tag with a specific key, regardless of
	//    the tag value.
	//
	//    * vpc-id - The ID of the VPC specified when the security group was created.
	Filters []Filter `locationName:"Filter" locationNameList:"Filter" type:"list"`

	// The IDs of the security groups. Required for security groups in a nondefault
	// VPC.
	//
	// Default: Describes all your security groups.
	GroupIds []string `locationName:"GroupId" locationNameList:"groupId" type:"list"`

	// [EC2-Classic and default VPC only] The names of the security groups. You
	// can specify either the security group name or the security group ID. For
	// security groups in a nondefault VPC, use the group-name filter to describe
	// security groups by name.
	//
	// Default: Describes all your security groups.
	GroupNames []string `locationName:"GroupName" locationNameList:"GroupName" type:"list"`

	// The maximum number of results to return in a single call. To retrieve the
	// remaining results, make another request with the returned NextToken value.
	// This value can be between 5 and 1000. If this parameter is not specified,
	// then all results are returned.
	MaxResults *int64 `min:"5" type:"integer"`

	// The token to request the next page of results.
	NextToken *string `type:"string"`
}

// String returns the string representation
func (s DescribeSecurityGroupsInput) String() string {
	return awsutil.Prettify(s)
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DescribeSecurityGroupsInput) Validate() error {
	invalidParams := aws.ErrInvalidParams{Context: "DescribeSecurityGroupsInput"}
	if s.MaxResults != nil && *s.MaxResults < 5 {
		invalidParams.Add(aws.NewErrParamMinValue("MaxResults", 5))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DescribeSecurityGroupsResult
type DescribeSecurityGroupsOutput struct {
	_ struct{} `type:"structure"`

	// The token to use to retrieve the next page of results. This value is null
	// when there are no more results to return.
	NextToken *string `locationName:"nextToken" type:"string"`

	// Information about the security groups.
	SecurityGroups []SecurityGroup `locationName:"securityGroupInfo" locationNameList:"item" type:"list"`
}

// String returns the string representation
func (s DescribeSecurityGroupsOutput) String() string {
	return awsutil.Prettify(s)
}

const opDescribeSecurityGroups = "DescribeSecurityGroups"

// DescribeSecurityGroupsRequest returns a request value for making API operation for
// Amazon Elastic Compute Cloud.
//
// Describes the specified security groups or all of your security groups.
//
// A security group is for use with instances either in the EC2-Classic platform
// or in a specific VPC. For more information, see Amazon EC2 Security Groups
// (https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/using-network-security.html)
// in the Amazon Elastic Compute Cloud User Guide and Security Groups for Your
// VPC (https://docs.aws.amazon.com/AmazonVPC/latest/UserGuide/VPC_SecurityGroups.html)
// in the Amazon Virtual Private Cloud User Guide.
//
//    // Example sending a request using DescribeSecurityGroupsRequest.
//    req := client.DescribeSecurityGroupsRequest(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
//
// Please also see https://docs.aws.amazon.com/goto/WebAPI/ec2-2016-11-15/DescribeSecurityGroups
func (c *Client) DescribeSecurityGroupsRequest(input *DescribeSecurityGroupsInput) DescribeSecurityGroupsRequest {
	op := &aws.Operation{
		Name:       opDescribeSecurityGroups,
		HTTPMethod: "POST",
		HTTPPath:   "/",
		Paginator: &aws.Paginator{
			InputTokens:     []string{"NextToken"},
			OutputTokens:    []string{"NextToken"},
			LimitToken:      "MaxResults",
			TruncationToken: "",
		},
	}

	if input == nil {
		input = &DescribeSecurityGroupsInput{}
	}

	req := c.newRequest(op, input, &DescribeSecurityGroupsOutput{})
	return DescribeSecurityGroupsRequest{Request: req, Input: input, Copy: c.DescribeSecurityGroupsRequest}
}

// DescribeSecurityGroupsRequest is the request type for the
// DescribeSecurityGroups API operation.
type DescribeSecurityGroupsRequest struct {
	*aws.Request
	Input *DescribeSecurityGroupsInput
	Copy  func(*DescribeSecurityGroupsInput) DescribeSecurityGroupsRequest
}

// Send marshals and sends the DescribeSecurityGroups API request.
func (r DescribeSecurityGroupsRequest) Send(ctx context.Context) (*DescribeSecurityGroupsResponse, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &DescribeSecurityGroupsResponse{
		DescribeSecurityGroupsOutput: r.Request.Data.(*DescribeSecurityGroupsOutput),
		response:                     &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// NewDescribeSecurityGroupsRequestPaginator returns a paginator for DescribeSecurityGroups.
// Use Next method to get the next page, and CurrentPage to get the current
// response page from the paginator. Next will return false, if there are
// no more pages, or an error was encountered.
//
// Note: This operation can generate multiple requests to a service.
//
//   // Example iterating over pages.
//   req := client.DescribeSecurityGroupsRequest(input)
//   p := ec2.NewDescribeSecurityGroupsRequestPaginator(req)
//
//   for p.Next(context.TODO()) {
//       page := p.CurrentPage()
//   }
//
//   if err := p.Err(); err != nil {
//       return err
//   }
//
func NewDescribeSecurityGroupsPaginator(req DescribeSecurityGroupsRequest) DescribeSecurityGroupsPaginator {
	return DescribeSecurityGroupsPaginator{
		Pager: aws.Pager{
			NewRequest: func(ctx context.Context) (*aws.Request, error) {
				var inCpy *DescribeSecurityGroupsInput
				if req.Input != nil {
					tmp := *req.Input
					inCpy = &tmp
				}

				newReq := req.Copy(inCpy)
				newReq.SetContext(ctx)
				return newReq.Request, nil
			},
		},
	}
}

// DescribeSecurityGroupsPaginator is used to paginate the request. This can be done by
// calling Next and CurrentPage.
type DescribeSecurityGroupsPaginator struct {
	aws.Pager
}

func (p *DescribeSecurityGroupsPaginator) CurrentPage() *DescribeSecurityGroupsOutput {
	return p.Pager.CurrentPage().(*DescribeSecurityGroupsOutput)
}

// DescribeSecurityGroupsResponse is the response type for the
// DescribeSecurityGroups API operation.
type DescribeSecurityGroupsResponse struct {
	*DescribeSecurityGroupsOutput

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// DescribeSecurityGroups request.
func (r *DescribeSecurityGroupsResponse) SDKResponseMetdata() *aws.Response {
	return r.response
}
