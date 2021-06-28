# TEMP TODO FILE

# Mockup for Recordings and Scheduling

## Move some things around

- [ ] move social icons under Profile image
- [ ] move External Actions to top right, below video element
- [ ] disable/hide chat panel + chat icon when there is no Recording, nor Live video playing
- [ ] add tab bar below tags list

### Add more local react States
- add offline / no-video state
- add recordings[] when comes in from config
- add schedule[] when comes in from config
- add route tracker
- add tab tracker

## Add Tab bar
Tab bar includes:
- `About` - User custom info
- `Videos` - display if user has Recordings
- `Schedule` - display tab if schedule info exists


## Routing, Url Handling

#### Recording urls
`server.com/recordings`
`server.com/recordings/id123`


#### Schedule urls
`server.com/schedule`
`server.com/schedule/id123`

### Todo
- [ ] modify server side go to just load up index.html/app.js when url routes to /recording or /schedule
- [ ] update app js to detect url route and display appropriate tab content


## Recordings

### `server.com/recordings`
- [ ] don't show chat elements
- [ ] list avilable recordings, display list similar to directory.owncast.
- [ ] display recording length
- [ ] display num views?


### `server.com/recordings/id123`
- [ ] display video, full size with recording loaded
- [ ] display chat
- [ ] do not enable chat message input
- [ ] render chat messages as they came in relative to video timestamp


## Schedule
- [ ] don't show chat elements

### `server.com/schedule`
- [ ] list items ASC


### `server.com/schedule/id123`
- [ ] display info


