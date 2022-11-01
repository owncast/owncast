/*
The Owncast Syncer.

It will try to sync the play time between clients.

Great for an LAN party.
*/

import { URL_PLAY_INSTRUCTIONS } from "../utils/constants.js";
import { getCurrentlyPlayingSegment } from "../utils/helpers.js";

export default class Syncer{
    constructor(player){
        this.player = player;
        this.clockSkewMs = 0;
        console.log("Syncer created")

        setInterval(() => {
            this.sync();
        }, 2000);
    }

    setClockSkew(ms) {
        this.clockSkewMs = ms;
    }

    async sync(){
        const options={
            method: 'GET'
        };
        const rawResponse = await fetch(URL_PLAY_INSTRUCTIONS,options);
        const content = await rawResponse.json();
        const expectedLatency = content.ExpectedLatency;

        const tech = this.player.tech({ IWillNotUseThisInPlugins: true });
        try {
            const segment = getCurrentlyPlayingSegment(tech);
            if (!segment || !segment.dateTimeObject) {
              return;
            }
            const intoSegment = segment.duration - (segment.end - this.player.currentTime());
            const segmentStartTime = segment.dateTimeObject.getTime();
            const segmentTime = segmentStartTime + intoSegment * 1000;
            const now = new Date().getTime() + this.clockSkewMs; //Server now
            const latency = (now - segmentTime) / 1000;

            // console.log("Current latency is",latency,", expected latency is",expectedLatency)

            if (Math.abs(latency-expectedLatency)<0.1 || expectedLatency==0){
                console.log("Current latency is",latency,", expected latency is",expectedLatency,", no need to sync")
                return;
            }else{
                let targetPlayTime = this.player.currentTime() + (latency - expectedLatency)
                if (targetPlayTime < 0) {
                    targetPlayTime = 0;
                }
                // if (targetPlayTime > segment.end) {
                //     targetPlayTime = segment.end;
                // }

                console.log("Current latency is",latency,", expected latency is",expectedLatency,", move ",targetPlayTime-this.player.currentTime(),"forward")
                this.player.currentTime(targetPlayTime);
            }
          } catch (err) {
            console.warn(err);
          }

    }
}