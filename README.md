# yearofcommits

Systray program in Go to show how many days in a row I contributed something on GitHub.

```
dep ensure
go install
yearofcommits -u github_user -t github_api_token
```

If you're using OSX, you can add this program to `launchd`. Create `/Library/LaunchAgents/yearofcommits.plist` file with the following XML:

```
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>yearofcommits</string>

    <key>OnDemand</key>
    <false/>

    <key>UserName</key>
    <string>MAC_USER</string>

    <key>GroupName</key>
    <string>MAC_GROUP</string>

    <key>ProgramArguments</key>
    <array>
            <string>/go/bin/yearofcommits</string>
            <string>-u</string>
            <string>GITHUB_USER</string>
            <string>-t</string>
            <string>GITHUB_TOKEN</string>
    </array>
</dict>
</plist>
```

Change `/go/bin/yearofcommits` to be in your $GOBIN.

Run:
```
sudo launchctl load /Library/LaunchAgents/yearofcommits.plist
```