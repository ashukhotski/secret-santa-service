# secret-santa-service
Secret Santa Service is a bot that helps to organise Secret Santa event in your company via Slack and share the joy of winter holidays.

API is written in Go. There are 4 REST commands, all of them are POST:
1. /get accepts Slack channel details, user ID, year and returns user's Secret Santa match and their postal address if Secret Santa party exists, or returns error otherwise.
2. /initialize accepts Slack channel details, user ID, user's Slack response URL, user's postal address and creates a new Secret Santa party for a given Slack channel and sets the user as a host, so that no one else can manage the event (i.e., randomize pairs before there are enough people registered to participate in the event). Returns error if the party has already been created for a given Slack channel by the user with a different ID.
3. /participate acccepts Slack channel details, user ID, user's Slack response URL, user's postal address and adds the user to the party initialized by /initialize command if it exists, or returns error otherwise.
4. /randomize accepts Slack channel details, user ID and, if the party exists and the user executing the command is the host, randomly creates pairs, followed a call to each registered participant's Slack response URL to notify participans about their matches in the Slack channel where the party is hosted, or returns error otherwise.

Microservice is containerized and can be built using docker-compose command.

Data is stored in MongoDB.

Once the app is built, it should be integrated with Slack. Refer to Slack app creation and management documentation for this purpose. 

Make the upcoming Christmas and New Year Eve special! Ho! Ho! Ho!
