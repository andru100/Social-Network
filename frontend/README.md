# Dating App Frontend

React based front end to Dating app project. This is in the final phases of being changed to a full Social Network. To be used in conjunction with "Dating Api" a go based backend.

The project is a work in progress. Currently landing page will defualt to sign in page if user doesn't have a valid JWT session. User can click signup if they haven't already. It will go to users home page if they are signed in. On the homepage users can post comments, view there photos and view a news feed containing all users comments and reply to them.

Users are able to reply to peoples posts and like them.

Currently, on profile setup page clicking the profile pic allows you to add your own. You can also edit you bio and add photos to your album which can be viewed on your profile by clicking the media tab. Cliking go to profile button takes users to there profile where you should be able to see your files and post comments. You can like a post by clicking the heart icon.

# To Do

SSO sign on.

2MFA for password reset.

Edit the show likes/replys function so that it shows only that posts likes.e 

Improving likes feature so a user can only like a post once and has ability to unlike.

Ability to delete images / posts and change username.

Create search and # tag feature. 

Add friends feature.

Improve UI / CSS.  

# To run

Currently all requests to the backend are set to localhost. If deploying on a seperate machine, change all mentions of local host to your servers ip address.

npm start