#!/usr/bin/env node

/**
 * ocTestStream.js - RTMP test stream utility
 * 
 * Requirements:
 *   - Node.js
 *   - ffmpeg (a recent version with loop video support)
 *   - a Sans family font (for overlay text)
 * 
 * Example: node ocTestStream.js ~/Downloads/*.mp4 rtmp://127.0.0.1/live/abc123
 */

const fs = require('fs');
const path = require('path');
const { spawn, exec } = require('child_process');
const os = require('os');

class OcTestStream {
    constructor() {
        this.ffmpegExec = null;
        this.destinationHost = 'rtmp://127.0.0.1/live/abc123';
        this.videoFiles = [];
        this.listFile = 'list.txt';
    }

    /**
     * Find ffmpeg executable across different platforms
     */
    async findFFmpeg() {
        const ffmpegExecs = ['ffmpeg', 'ffmpeg.exe'];
        const ffmpegPaths = ['./', '../', ''];

        // Try local paths first
        for (const execName of ffmpegExecs) {
            for (const execPath of ffmpegPaths) {
                const fullPath = path.join(execPath, execName);
                try {
                    await this.checkExecutable(fullPath);
                    this.ffmpegExec = fullPath;
                    return true;
                } catch (e) {
                    // Continue searching
                }
            }
        }

        // Try system PATH
        for (const execName of ffmpegExecs) {
            try {
                await this.checkExecutable(execName);
                this.ffmpegExec = execName;
                return true;
            } catch (e) {
                // Continue searching
            }
        }

        return false;
    }

    /**
     * Check if an executable exists and is accessible
     */
    checkExecutable(execPath) {
        return new Promise((resolve, reject) => {
            exec(`"${execPath}" -version`, (error, stdout, stderr) => {
                if (error) {
                    reject(error);
                } else {
                    resolve(stdout);
                }
            });
        });
    }

    /**
     * Get ffmpeg version information
     */
    async getFFmpegVersion() {
        try {
            const output = await this.checkExecutable(this.ffmpegExec);
            const versionMatch = output.match(/ffmpeg version ([^\s]+)/);
            return versionMatch ? versionMatch[1] : 'unknown';
        } catch (e) {
            return 'unknown';
        }
    }

    /**
     * Get ffmpeg executable path
     */
    async getFFmpegPath() {
        return new Promise((resolve) => {
            const command = os.platform() === 'win32' ? `where "${this.ffmpegExec}"` : `which "${this.ffmpegExec}"`;
            exec(command, (error, stdout, stderr) => {
                if (error) {
                    resolve('unknown');
                } else {
                    resolve(stdout.trim().split('\n')[0]);
                }
            });
        });
    }

    /**
     * Parse command line arguments
     */
    parseArguments(args) {
        if (args.includes('--help')) {
            this.showHelp();
            return false;
        }

        // Check if last argument is RTMP URL
        const lastArg = args[args.length - 1];
        if (lastArg && lastArg.includes('rtmp://')) {
            console.log('RTMP server is specified');
            this.destinationHost = lastArg;
            this.videoFiles = args.slice(0, -1);
        } else {
            console.log('RTMP server is not specified');
            this.videoFiles = args;
        }

        return true;
    }

    /**
     * Show help information
     */
    showHelp() {
        console.log('ocTestStream is used for sending pre-recorded or internal test content to an RTMP server.');
        console.log('Usage: node ocTestStream.js [VIDEO_FILES] [RTMP_DESTINATION]');
        console.log('VIDEO_FILES: path to one or multiple videos for sending to the RTMP server (optional)');
        console.log('RTMP_DESTINATION: URL of RTMP server with key (optional; default: rtmp://127.0.0.1/live/abc123)');
    }

    /**
     * Get system font path for overlay text
     */
    getSystemFont() {
        const platform = os.platform();
        
        if (platform === 'win32') {
            // Workaround lack of font config on Windows.
            const windir = process.env.WINDIR || 'C:\\Windows';
            return path.join(windir, 'fonts', 'arial.ttf').replace(/\\/g, '/').replace(/:/g, '\\\\:');
        } else {
            // Use default font.
            return '';
        }
    }

    /**
     * Stream internal test video pattern
     */
    streamTestPattern() {
        console.log(`Streaming internal test video loop to ${this.destinationHost}`);
        console.log('...press ctrl+c to exit');

        const font = this.getSystemFont();
        const fontParam = font ? `:fontfile=${font}` : '';

        const args = [
            '-hide_banner', '-loglevel', 'panic', '-nostdin', '-re', '-f', 'lavfi',
            '-i', 'testsrc=size=1280x720:rate=60[out0];sine=frequency=400:sample_rate=48000[out1]',
            '-vf', `[in]drawtext=fontsize=96:box=1:boxcolor=black@0.75:boxborderw=5${fontParam}:fontcolor=white:x=(w-text_w)/2:y=((h-text_h)/2)+((h-text_h)/-2):text='Owncast Test Stream',drawtext=fontsize=96:box=1:boxcolor=black@0.75:boxborderw=5${fontParam}:fontcolor=white:x=(w-text_w)/2:y=((h-text_h)/2)+((h-text_h)/2):text='%{gmtime\\:%H\\\\\\:%M\\\\\\:%S} UTC'[out]`,
            '-nal-hrd', 'cbr',
            '-metadata:s:v', 'encoder=test',
            '-vcodec', 'libx264',
            '-acodec', 'aac',
            '-preset', 'veryfast',
            '-profile:v', 'baseline',
            '-tune', 'zerolatency',
            '-bf', '0',
            '-g', '0',
            '-b:v', '6320k',
            '-b:a', '160k',
            '-ac', '2',
            '-ar', '48000',
            '-minrate', '6320k',
            '-maxrate', '6320k',
            '-bufsize', '6320k',
            '-muxrate', '6320k',
            '-r', '60',
            '-pix_fmt', 'yuv420p',
            '-color_range', '1', '-colorspace', '1', '-color_primaries', '1', '-color_trc', '1',
            '-flags:v', '+global_header',
            '-bsf:v', 'dump_extra',
            '-x264-params', 'nal-hrd=cbr:min-keyint=2:keyint=2:scenecut=0:bframes=0',
            '-f', 'flv', this.destinationHost
        ];

        this.streamWithArgs(args);
    }

    /**
     * Stream with arguments
     */
    streamWithArgs(args) {
        const ffmpeg = spawn(this.ffmpegExec, args, { detached: false, stdio: 'inherit' });

        ffmpeg.on('error', (err) => {
            console.error('Error starting ffmpeg:', err);
            process.exit(1);
        });

        process.on('SIGINT', () => {
            if (os.platform() === 'win32') {
                // Because ffmpeg can spawn children, on windows we need to use taskkill.
                // Using ffmpeg.kill on Windows may result in orphaned processes.
                exec(`taskkill /T /PID ${ffmpeg.pid}`);
            } else {
                ffmpeg.kill('SIGINT');
            }
        });

        process.on('exit', () => {
            if (os.platform() === 'win32') {
                exec(`taskkill /F /T /PID ${ffmpeg.pid}`);
            } else {
                ffmpeg.kill('SIGINT');
            }
        });
    }

    /**
     * Create playlist file for video files
     */
    createPlaylist() {
        try {
            const content = this.videoFiles.map(file => `file '${file}'`).join('\n');
            fs.writeFileSync(this.listFile, content);
            return true;
        } catch (err) {
            console.error('Error creating playlist file:', err);
            return false;
        }
    }

    /**
     * Validate video files exist
     */
    validateFiles() {
        for (const file of this.videoFiles) {
            if (!fs.existsSync(file)) {
                console.error(`ERROR: File not found: ${file}`);
                return false;
            }
        }
        return true;
    }

    /**
     * Stream video files
     */
    streamVideoFiles() {
        if (!this.validateFiles()) {
            process.exit(1);
        }

        if (!this.createPlaylist()) {
            process.exit(1);
        }

        console.log(`Streaming a loop of ${this.videoFiles.length} video(s) to ${this.destinationHost}`);
        if (this.videoFiles.length > 1) {
            console.log('Warning: If these files differ greatly in formats, transitioning from one to another may not always work correctly.');
        }
        console.log(this.videoFiles.join(' '));
        console.log('...press ctrl+c to exit');

        const font = this.getSystemFont();
        const fontParam = font ? `:fontfile=${font}` : '';

        const args = [
            '-hide_banner', '-loglevel', 'panic', '-nostdin', '-stream_loop', '-1', '-re', '-f', 'concat',
            '-safe', '0',
            '-i', this.listFile,
            '-vcodec', 'libx264',
            '-profile:v', 'high',
            '-g', '48',
            '-r', '24',
            '-sc_threshold', '0',
            '-b:v', '1300k',
            '-preset', 'veryfast',
            '-acodec', 'copy',
            '-vf', `drawtext=fontsize=96:box=1:boxcolor=black@0.75:boxborderw=5${fontParam}:fontcolor=white:x=(w-text_w)/2:y=((h-text_h)/2)+((h-text_h)/4):text='%{gmtime\\:%H\\\\\\:%M\\\\\\:%S}'`,
            '-f', 'flv', this.destinationHost
        ];

        this.streamWithArgs(args);
    }

    /**
     * Cleanup temporary files
     */
    cleanup() {
        try {
            if (fs.existsSync(this.listFile)) {
                fs.unlinkSync(this.listFile);
            }
        } catch (err) {
            console.error('Error during cleanup:', err);
        }
    }

    /**
     * Main execution function
     */
    async run() {
        // Parse command line arguments (skip node and script name)
        const args = process.argv.slice(2);
        
        if (!this.parseArguments(args)) {
            return;
        }

        // Find ffmpeg executable
        if (!(await this.findFFmpeg())) {
            console.error('ERROR: ffmpeg was not found in path or in the current directory! Please install ffmpeg before using this script.');
            process.exit(1);
        }

        // Show ffmpeg information
        const version = await this.getFFmpegVersion();
        const execPath = await this.getFFmpegPath();
        console.log(`ffmpeg executable: ${this.ffmpegExec} (${version})`);
        console.log(`ffmpeg path: ${execPath}`);

        // Stream based on whether files were provided
        if (this.videoFiles.length === 0) {
            this.streamTestPattern();
        } else {
            this.streamVideoFiles();
        }
    }
}

// Run if called directly
if (require.main === module) {
    const stream = new OcTestStream();
    stream.run().catch(err => {
        console.error('Error:', err);
        process.exit(1);
    });
}

module.exports = OcTestStream;
