# What is this

This is the backend part of my project called "chanl".\
See the frontend [here](https://github.com/kutoru/chanl-frontend).\
The project is generally supposed to be a text-only Discord clone.\
With this project I wanted to learn WebSockets and fullstack web development.

# Project status

Just started

# The idea

The general idea is that the entire app consists of channels.\
Different channels kinda have different purposes and permissions.\
For each channel, there is at least one user that can write in it.\
Channels have parent-child relations between eachother.\
User can access child channels from their parent channels.

These are all the channels and their relations:
- Global -> Private -> Server -> Room
- Global -> Personal -> Friend

So here, for an example, the Global channel doesn't have any parent channels and has two children, Private and Personal channels. The Private channel is Server's parent, and Room is Server's child.

Here are descriptions and permissions of each channel:
- Global: the main channel.\
Contains Private and Personal channels.\
Available to all users.\
Only the developers can write in the channel.
- Private: unique to each user.\
Contains a list of servers that the user has joined.\
Available to the members of servers that the user has joined.\
Only the user can write in the channel.
- Server: can be created by users.\
Contains a list of rooms that the owner has created.\
Available to the members of the server.\
Only the owner can write in the channel.
- Room: can be created by the owner of the Server.\
Doesn't have any children.\
Available to the members of the server.\
All members can write in the channel.
- Personal: unique to each user.\
Contains a list of user's friends.\
Available to the user's friends.\
Only the user can write in the channel.
- Friend: basically a dm.\
Doesn't have any children.\
Available to the two people who are in the dm.\
Both can write in the channel.

Not sure how practical the entire idea is but I think it will be fun to implement
