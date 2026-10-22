let me tell you the entire intended flow then you criticise the code that is written and assume the code that is not written to be according to what i said

there are 4 stages of initialisations when program starts,

1. the config is loaded in memory, it will include some general behaviour info then list out servers witht their own rules

2. we create a channel that is monitored by the workers that are yet to be created will need the structure of this channel to also accomodate the writer, as it will be needed when we are sending the response back; i'm thinking we create this struct called job, it will have the raw request, the writer and some id so we can track it better.

3. then we start the http listener that just listens to whatever is coming at the proxy server, and puts it in the channel. the max concurrent connections is limited by using the semaphore structure. the listeners here has to keep the connection open till the backend sends the reponse otherwise the client will get the wrong signal about the connection being closed.

4. then the workers are initiated it picks one from the channel, asks the algorithm where to send this to and what to change about the raw request, the algorithm checks if where this is coming from, is it from the backend servers that are listed in config or not and tells exactly what to do. the workers then make their own concurrent requests (each a go function, i'm thinking of capping it at 10) and whenever they are done, go back to the same loop while the main program stops in the blocking condition so nothing stops working
