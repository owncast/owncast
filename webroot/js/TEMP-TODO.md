# TEMP TODO FILE

# ~~Mockup for Recordings and Scheduling~~

# Mockup for Fediverse Social, Tabbed User Content

This used to be setting up to display Recordings, but the progress can be used towards Fediverse work.

- Rearranges some logic around when to display the chat panel vs when video is playing
- Improves user content styling for improved positioning across screen sizes.
- Add accessible Tab Bar navigation

## Move some things around

- [x] move social icons under Profile image
- [x] move External Actions to top right, below video element
- [x] disable/hide chat panel + chat icon when there is no Recording, nor Live video playing
- [ ] add tab bar below tags list
- [ ] style Follow on Fediverse Modal

### Add more local React States

- [ ] add offline / no-video state (? what was this again?)
- [ ] add tab states
- [ ] **DEFER** add route states
- [ ] **DEFER** add recordings[] when comes in from config
- [ ] **DEFER** add schedule[] when comes in from config

## Add Tab bar

Tab bar includes:

- `About` - User custom info
- `Followers` - display tab if schedule info exists
- **DEFER** `Videos` - display if user has Recordings
- **DEFER** `Schedule` - display tab if schedule info exists

## **DEFER?** Routing, Url Handling

- do we need it for Followers?

  #### Recording urls

  - `server.com/recordings`
  - `server.com/recordings/id123`

  #### Schedule urls

  - `server.com/schedule`
  - `server.com/schedule/id123`

  #### Followers Url?

  ### Todo

  - [ ] modify server side go to just load up index.html/app.js when url routes to /recording or /schedule
  - [ ] update app js to detect url route and display appropriate tab content

## **DEFER** Recordings

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
