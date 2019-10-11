// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package cloudformation

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awsutil"
)

// The input for CreateStack action.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/cloudformation-2010-05-15/CreateStackInput
type CreateStackInput struct {
	_ struct{} `type:"structure"`

	// In some cases, you must explicity acknowledge that your stack template contains
	// certain capabilities in order for AWS CloudFormation to create the stack.
	//
	//    * CAPABILITY_IAM and CAPABILITY_NAMED_IAM Some stack templates might include
	//    resources that can affect permissions in your AWS account; for example,
	//    by creating new AWS Identity and Access Management (IAM) users. For those
	//    stacks, you must explicitly acknowledge this by specifying one of these
	//    capabilities. The following IAM resources require you to specify either
	//    the CAPABILITY_IAM or CAPABILITY_NAMED_IAM capability. If you have IAM
	//    resources, you can specify either capability. If you have IAM resources
	//    with custom names, you must specify CAPABILITY_NAMED_IAM. If you don't
	//    specify either of these capabilities, AWS CloudFormation returns an InsufficientCapabilities
	//    error. If your stack template contains these resources, we recommend that
	//    you review all permissions associated with them and edit their permissions
	//    if necessary. AWS::IAM::AccessKey (https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iam-accesskey.html)
	//    AWS::IAM::Group (https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iam-group.html)
	//    AWS::IAM::InstanceProfile (https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iam-instanceprofile.html)
	//    AWS::IAM::Policy (https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iam-policy.html)
	//    AWS::IAM::Role (https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iam-role.html)
	//    AWS::IAM::User (https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iam-user.html)
	//    AWS::IAM::UserToGroupAddition (https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iam-addusertogroup.html)
	//    For more information, see Acknowledging IAM Resources in AWS CloudFormation
	//    Templates (http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-iam-template.html#capabilities).
	//
	//    * CAPABILITY_AUTO_EXPAND Some template contain macros. Macros perform
	//    custom processing on templates; this can include simple actions like find-and-replace
	//    operations, all the way to extensive transformations of entire templates.
	//    Because of this, users typically create a change set from the processed
	//    template, so that they can review the changes resulting from the macros
	//    before actually creating the stack. If your stack template contains one
	//    or more macros, and you choose to create a stack directly from the processed
	//    template, without first reviewing the resulting changes in a change set,
	//    you must acknowledge this capability. This includes the AWS::Include (http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/create-reusable-transform-function-snippets-and-add-to-your-template-with-aws-include-transform.html)
	//    and AWS::Serverless (http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/transform-aws-serverless.html)
	//    transforms, which are macros hosted by AWS CloudFormation. Change sets
	//    do not currently support nested stacks. If you want to create a stack
	//    from a stack template that contains macros and nested stacks, you must
	//    create the stack directly from the template using this capability. You
	//    should only create stacks directly from a stack template that contains
	//    macros if you know what processing the macro performs. Each macro relies
	//    on an underlying Lambda service function for processing stack templates.
	//    Be aware that the Lambda function owner can update the function operation
	//    without AWS CloudFormation being notified. For more information, see Using
	//    AWS CloudFormation Macros to Perform Custom Processing on Templates (http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/template-macros.html).
	Capabilities []Capability `type:"list"`

	// A unique identifier for this CreateStack request. Specify this token if you
	// plan to retry requests so that AWS CloudFormation knows that you're not attempting
	// to create a stack with the same name. You might retry CreateStack requests
	// to ensure that AWS CloudFormation successfully received them.
	//
	// All events triggered by a given stack operation are assigned the same client
	// request token, which you can use to track operations. For example, if you
	// execute a CreateStack operation with the token token1, then all the StackEvents
	// generated by that operation will have ClientRequestToken set as token1.
	//
	// In the console, stack operations display the client request token on the
	// Events tab. Stack operations that are initiated from the console use the
	// token format Console-StackOperation-ID, which helps you easily identify the
	// stack operation . For example, if you create a stack using the console, each
	// stack event would be assigned the same token in the following format: Console-CreateStack-7f59c3cf-00d2-40c7-b2ff-e75db0987002.
	ClientRequestToken *string `min:"1" type:"string"`

	// Set to true to disable rollback of the stack if stack creation failed. You
	// can specify either DisableRollback or OnFailure, but not both.
	//
	// Default: false
	DisableRollback *bool `type:"boolean"`

	// Whether to enable termination protection on the specified stack. If a user
	// attempts to delete a stack with termination protection enabled, the operation
	// fails and the stack remains unchanged. For more information, see Protecting
	// a Stack From Being Deleted (http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-protect-stacks.html)
	// in the AWS CloudFormation User Guide. Termination protection is disabled
	// on stacks by default.
	//
	// For nested stacks (http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-nested-stacks.html),
	// termination protection is set on the root stack and cannot be changed directly
	// on the nested stack.
	EnableTerminationProtection *bool `type:"boolean"`

	// The Simple Notification Service (SNS) topic ARNs to publish stack related
	// events. You can find your SNS topic ARNs using the SNS console or your Command
	// Line Interface (CLI).
	NotificationARNs []string `type:"list"`

	// Determines what action will be taken if stack creation fails. This must be
	// one of: DO_NOTHING, ROLLBACK, or DELETE. You can specify either OnFailure
	// or DisableRollback, but not both.
	//
	// Default: ROLLBACK
	OnFailure OnFailure `type:"string" enum:"true"`

	// A list of Parameter structures that specify input parameters for the stack.
	// For more information, see the Parameter (https://docs.aws.amazon.com/AWSCloudFormation/latest/APIReference/API_Parameter.html)
	// data type.
	Parameters []Parameter `type:"list"`

	// The template resource types that you have permissions to work with for this
	// create stack action, such as AWS::EC2::Instance, AWS::EC2::*, or Custom::MyCustomInstance.
	// Use the following syntax to describe template resource types: AWS::* (for
	// all AWS resource), Custom::* (for all custom resources), Custom::logical_ID
	// (for a specific custom resource), AWS::service_name::* (for all resources
	// of a particular AWS service), and AWS::service_name::resource_logical_ID
	// (for a specific AWS resource).
	//
	// If the list of resource types doesn't include a resource that you're creating,
	// the stack creation fails. By default, AWS CloudFormation grants permissions
	// to all resource types. AWS Identity and Access Management (IAM) uses this
	// parameter for AWS CloudFormation-specific condition keys in IAM policies.
	// For more information, see Controlling Access with AWS Identity and Access
	// Management (https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-iam-template.html).
	ResourceTypes []string `type:"list"`

	// The Amazon Resource Name (ARN) of an AWS Identity and Access Management (IAM)
	// role that AWS CloudFormation assumes to create the stack. AWS CloudFormation
	// uses the role's credentials to make calls on your behalf. AWS CloudFormation
	// always uses this role for all future operations on the stack. As long as
	// users have permission to operate on the stack, AWS CloudFormation uses this
	// role even if the users don't have permission to pass it. Ensure that the
	// role grants least privilege.
	//
	// If you don't specify a value, AWS CloudFormation uses the role that was previously
	// associated with the stack. If no role is available, AWS CloudFormation uses
	// a temporary session that is generated from your user credentials.
	RoleARN *string `min:"20" type:"string"`

	// The rollback triggers for AWS CloudFormation to monitor during stack creation
	// and updating operations, and for the specified monitoring period afterwards.
	RollbackConfiguration *RollbackConfiguration `type:"structure"`

	// The name that is associated with the stack. The name must be unique in the
	// region in which you are creating the stack.
	//
	// A stack name can contain only alphanumeric characters (case sensitive) and
	// hyphens. It must start with an alphabetic character and cannot be longer
	// than 128 characters.
	//
	// StackName is a required field
	StackName *string `type:"string" required:"true"`

	// Structure containing the stack policy body. For more information, go to Prevent
	// Updates to Stack Resources (https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/protect-stack-resources.html)
	// in the AWS CloudFormation User Guide. You can specify either the StackPolicyBody
	// or the StackPolicyURL parameter, but not both.
	StackPolicyBody *string `min:"1" type:"string"`

	// Location of a file containing the stack policy. The URL must point to a policy
	// (maximum size: 16 KB) located in an S3 bucket in the same region as the stack.
	// You can specify either the StackPolicyBody or the StackPolicyURL parameter,
	// but not both.
	StackPolicyURL *string `min:"1" type:"string"`

	// Key-value pairs to associate with this stack. AWS CloudFormation also propagates
	// these tags to the resources created in the stack. A maximum number of 50
	// tags can be specified.
	Tags []Tag `type:"list"`

	// Structure containing the template body with a minimum length of 1 byte and
	// a maximum length of 51,200 bytes. For more information, go to Template Anatomy
	// (https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/template-anatomy.html)
	// in the AWS CloudFormation User Guide.
	//
	// Conditional: You must specify either the TemplateBody or the TemplateURL
	// parameter, but not both.
	TemplateBody *string `min:"1" type:"string"`

	// Location of file containing the template body. The URL must point to a template
	// (max size: 460,800 bytes) that is located in an Amazon S3 bucket. For more
	// information, go to the Template Anatomy (https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/template-anatomy.html)
	// in the AWS CloudFormation User Guide.
	//
	// Conditional: You must specify either the TemplateBody or the TemplateURL
	// parameter, but not both.
	TemplateURL *string `min:"1" type:"string"`

	// The amount of time that can pass before the stack status becomes CREATE_FAILED;
	// if DisableRollback is not set or is set to false, the stack will be rolled
	// back.
	TimeoutInMinutes *int64 `min:"1" type:"integer"`
}

// String returns the string representation
func (s CreateStackInput) String() string {
	return awsutil.Prettify(s)
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *CreateStackInput) Validate() error {
	invalidParams := aws.ErrInvalidParams{Context: "CreateStackInput"}
	if s.ClientRequestToken != nil && len(*s.ClientRequestToken) < 1 {
		invalidParams.Add(aws.NewErrParamMinLen("ClientRequestToken", 1))
	}
	if s.RoleARN != nil && len(*s.RoleARN) < 20 {
		invalidParams.Add(aws.NewErrParamMinLen("RoleARN", 20))
	}

	if s.StackName == nil {
		invalidParams.Add(aws.NewErrParamRequired("StackName"))
	}
	if s.StackPolicyBody != nil && len(*s.StackPolicyBody) < 1 {
		invalidParams.Add(aws.NewErrParamMinLen("StackPolicyBody", 1))
	}
	if s.StackPolicyURL != nil && len(*s.StackPolicyURL) < 1 {
		invalidParams.Add(aws.NewErrParamMinLen("StackPolicyURL", 1))
	}
	if s.TemplateBody != nil && len(*s.TemplateBody) < 1 {
		invalidParams.Add(aws.NewErrParamMinLen("TemplateBody", 1))
	}
	if s.TemplateURL != nil && len(*s.TemplateURL) < 1 {
		invalidParams.Add(aws.NewErrParamMinLen("TemplateURL", 1))
	}
	if s.TimeoutInMinutes != nil && *s.TimeoutInMinutes < 1 {
		invalidParams.Add(aws.NewErrParamMinValue("TimeoutInMinutes", 1))
	}
	if s.RollbackConfiguration != nil {
		if err := s.RollbackConfiguration.Validate(); err != nil {
			invalidParams.AddNested("RollbackConfiguration", err.(aws.ErrInvalidParams))
		}
	}
	if s.Tags != nil {
		for i, v := range s.Tags {
			if err := v.Validate(); err != nil {
				invalidParams.AddNested(fmt.Sprintf("%s[%v]", "Tags", i), err.(aws.ErrInvalidParams))
			}
		}
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// The output for a CreateStack action.
// Please also see https://docs.aws.amazon.com/goto/WebAPI/cloudformation-2010-05-15/CreateStackOutput
type CreateStackOutput struct {
	_ struct{} `type:"structure"`

	// Unique identifier of the stack.
	StackId *string `type:"string"`
}

// String returns the string representation
func (s CreateStackOutput) String() string {
	return awsutil.Prettify(s)
}

const opCreateStack = "CreateStack"

// CreateStackRequest returns a request value for making API operation for
// AWS CloudFormation.
//
// Creates a stack as specified in the template. After the call completes successfully,
// the stack creation starts. You can check the status of the stack via the
// DescribeStacks API.
//
//    // Example sending a request using CreateStackRequest.
//    req := client.CreateStackRequest(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
//
// Please also see https://docs.aws.amazon.com/goto/WebAPI/cloudformation-2010-05-15/CreateStack
func (c *Client) CreateStackRequest(input *CreateStackInput) CreateStackRequest {
	op := &aws.Operation{
		Name:       opCreateStack,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &CreateStackInput{}
	}

	req := c.newRequest(op, input, &CreateStackOutput{})
	return CreateStackRequest{Request: req, Input: input, Copy: c.CreateStackRequest}
}

// CreateStackRequest is the request type for the
// CreateStack API operation.
type CreateStackRequest struct {
	*aws.Request
	Input *CreateStackInput
	Copy  func(*CreateStackInput) CreateStackRequest
}

// Send marshals and sends the CreateStack API request.
func (r CreateStackRequest) Send(ctx context.Context) (*CreateStackResponse, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &CreateStackResponse{
		CreateStackOutput: r.Request.Data.(*CreateStackOutput),
		response:          &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// CreateStackResponse is the response type for the
// CreateStack API operation.
type CreateStackResponse struct {
	*CreateStackOutput

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// CreateStack request.
func (r *CreateStackResponse) SDKResponseMetdata() *aws.Response {
	return r.response
}
