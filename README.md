# reserve-gw
## Reservation gateway improvement project

This project started out as a proof of concept to demonstrate possible improvements
to the online reservation systems of a couple of state parks. Two main goals were to
improve the search capabilty and decrease the time a patron spends obtaining results,
thus reducing his/her frustration level. The reservation systems were to work
transparently as they are and without impact to any development process or flow.

Note: these reservation systems consist of camping and day use sites across the two states,
Tennessee and Texas. During high seasons when lengths of stay at any one campsite is limited,
patrons try to book multiple campsites at a location to cover their dates of interest.
This currently is done only via a phone call to human reservation specialist.

### Goals:
#### Improve search capability:
+ Turn from relatively flat results based upon date ranges into multi-varied results of possible combinations over the date ranges;
+ Present the results in a more informative and sortable way to provide all the information available cleanly and quickly;

#### Decrease patron frustration:
+ Learn and maintain user preferences including favorite camping destinations and campsites;
+ Provide semi-autonomous search mechanism to find and notify user of locations, dates, and possible combinations;

### Status:
Currently, the project includes the interactive initial concept. The _first phase_ is to complete the backend result combinations construction for the frontend presentation. The _second phase_ would then tackle the semi-autonomous search mechanism, and a _third phase_ would implement a 
learning knowledge base for users.

### Building and running:
+ Clone the repository into your **Go src** directory
+ Change the working directory to the reserve-gw directory
+ Execute 'go build ./gateway'
+ Assuming it all builds correctly, execute './gateway'
+ Open a browser tab to localhost:8000/tnsp and begin to use the gateway

The prototype will start a server listening on the default port. The port value for the 
server is configured to a default value 8000, which can be overridden with a --port <value>
command line switch. The debug messages produced on the console do not yet have
(and may never have) a verbosity switch. If you change the port number, then browse to localhost:\<new port number\>/tnsp.
