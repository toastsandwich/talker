package main

/*
	main routine
		|
		|
		|	  server routine
		|			|			   routine
		|-----------|		   TCP         UDP
		|			|-----------|-----------|--------------------------------
		|			|			|			|								|
		|			|			|			|								|
		|			|			|			|								|
		|			.			.			.								|
		|			.			.			.								|
		|			.			.			.								|
		|						^			^								|
		|	  client routine	|			|								|
		|-----------|			|			|								|
		|			|			|  routines |				routine			|		 routine
		|			|		 TCP REQ	  UDP REQ			TCPRESP---------+--------UDP RESP
		|			|-----------|-----------|				|							|
		|			|			|			|				|							|
		|			|			|			|				|							|
		|			|			|			|				|							|
		.			.			.			.				|							|
		.			.			.			.				.							.
		.			.			.			.				.							.
		.			.			.			.				.							.
		.			.			.			.				.							.
*/

func main() {
	// set up a cli for getting local addr, remote addr, rport, lport, timeout
}
