This can be any text that makes sense to you.
```
[Unit]
Description=Owncast Service
```

This is where the "functional" parts of the service live.<br />
`WorkingDirectory` should be where you want the owncast folder to live.<br />
**Example:** ```WorkingDirectory=\var\www\owncast```

Similarly the `ExecStart` is the actual owncast binary.<br />
**Example:** ```ExecStart=\var\www\owncast\owncast```

```
[Service]
Type=simple
WorkingDirectory=[path_to_owncast_root_directory]
ExecStart=[path_to_owncast_executable]
Restart=on-failure
RestartSec=5
```

This just means, use runlevel 3 non-graphical.
```
[Install]
WantedBy=multi-user.target
```
