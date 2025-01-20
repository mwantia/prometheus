# Tools

* Worker receives new prompt
* Prompt is compared for certain keywords
  * Each internal tool can have certain keywords it listens to
  * If a keyword has been found (eg. `generate` or `random`) the tool assosiated with it gets included in the request
* Create new chat request with the available tools, previous messages and the new prompt
* Chat is being called and the worker waits for the final response
* If the response contains a tool call, the following applies:
  * The response is dissected to get the required name and arguments
  * The worker performs the internal call for the function, waits for the reply and writes it down in a complete sentence (eg. `The generated random number is '5'`)
  * This new sentence is added to a new chat request without any tools
  * Chat is being called and the worker waits for the final response
* The message content of the received response is validated and send to the user as response
