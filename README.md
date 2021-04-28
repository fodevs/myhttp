# myhttp
### This is a tool which makes http requests and prints the address of the request along with the MD5 hash of the response.

### Example
$> ./myhttp adjust.com

http://adjust.com d1b40e2a2ba488a054186e4ed0733f9752f66949


The tool is able to limit the number of parallel requests, to prevent exhausting local resources.
The tool may accept a flag '-parallel' to indicate this limit, default value is 10 if the flag is not provided.

$> ./myhttp -parallel 3 adjust.com google.com facebook.com yahoo.com