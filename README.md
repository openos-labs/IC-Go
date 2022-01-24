# IC-Go: Go Agent for the Internet Computer

Here is the beta version of go agent for IC

Please send us the issues you found, thanks!

You can find the examples to use IC-Go in agent_test.go

The implementations of IDL and principal are borrowed from [candid-go](https://github.com/aviate-labs/candid-go) and [principal](https://github.com/aviate-labs/principal-go), and fix the bugs

**[update 2022.1.24]** add the auto-decoder to decode the result to the data structure in Golang automatically, see utils/decode.go

	Tips for the auto-decoder:

    - Due to the limitation of member variable name in Golang, you should put the label for member variable in the tag of `ic:"**"`

	- You are supposed to define Enum class as a struct. The extra member variant with tag `ic:"EnumIndex"`, whose value is set as the final result's tag, is used to formulate the enumeration in such a "enum" struct   
	
	- "Option" should be defined as a struct, and the tag of the member variants in such a struct are "none" and "some", the member variant with tag "none" is initialized as 0, and once the option is nil, it will be set to 1  
	
	- "Tuple" is also defined as a struct, the tag for the first member in tuple should be `ic:"0"`, and so on.

    - "Nat" is also supposed to be defined as big.Int 
	
	- See more details in utils/decode.go and agent_test.go

    


**[update 2022.1.14]** add the new function of agent from pem file and show an example to decode the result to the data structure in Golang