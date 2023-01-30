#!/bin/bash

# This script will make API calls to the Owncast admin API and populate
# content for the page.

echo "Sending API calls to the Owncast admin to fake server config..."

# Server name

curl 'http://localhost:8080/api/admin/config/name' \
	-H 'Accept: */*' \
	-H 'Accept-Language: en-US,en;q=0.9' \
	-H 'Authorization: Basic YWRtaW46YWJjMTIz' \
	-H 'Cache-Control: no-cache' \
	-H 'Connection: keep-alive' \
	-H 'Content-Type: text/plain;charset=UTF-8' \
	-H 'Origin: http://localhost:8080' \
	-H 'Pragma: no-cache' \
	-H 'Referer: http://localhost:8080/admin/config-public-details/' \
	-H 'Sec-Fetch-Dest: empty' \
	-H 'Sec-Fetch-Mode: cors' \
	-H 'Sec-Fetch-Site: same-origin' \
	-H 'Sec-GPC: 1' \
	-H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36' \
	--data-raw '{"value":"My Cool Stream"}' \
	--compressed

# Description/Summary

curl 'http://localhost:8080/api/admin/config/serversummary' \
	-H 'Accept: */*' \
	-H 'Accept-Language: en-US,en;q=0.9' \
	-H 'Authorization: Basic YWRtaW46YWJjMTIz' \
	-H 'Cache-Control: no-cache' \
	-H 'Connection: keep-alive' \
	-H 'Content-Type: text/plain;charset=UTF-8' \
	-H 'Origin: http://localhost:8080' \
	-H 'Pragma: no-cache' \
	-H 'Referer: http://localhost:8080/admin/config-public-details/' \
	-H 'Sec-Fetch-Dest: empty' \
	-H 'Sec-Fetch-Mode: cors' \
	-H 'Sec-Fetch-Site: same-origin' \
	-H 'Sec-GPC: 1' \
	-H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36' \
	--data-raw '{"value":"Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. https://owncast.online"}' \
	--compressed

# Tags

curl 'http://localhost:8080/api/admin/config/tags' \
	-H 'Accept: */*' \
	-H 'Accept-Language: en-US,en;q=0.9' \
	-H 'Authorization: Basic YWRtaW46YWJjMTIz' \
	-H 'Cache-Control: no-cache' \
	-H 'Connection: keep-alive' \
	-H 'Content-Type: text/plain;charset=UTF-8' \
	-H 'Origin: http://localhost:8080' \
	-H 'Pragma: no-cache' \
	-H 'Referer: http://localhost:8080/admin/config-public-details/' \
	-H 'Sec-Fetch-Dest: empty' \
	-H 'Sec-Fetch-Mode: cors' \
	-H 'Sec-Fetch-Site: same-origin' \
	-H 'Sec-GPC: 1' \
	-H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36' \
	--data-raw '{"value":["owncast","streaming","testing","more","tags"]}' \
	--compressed

# Page content

curl 'http://localhost:8080/api/admin/config/pagecontent' \
	-H 'Accept: */*' \
	-H 'Accept-Language: en-US,en;q=0.9' \
	-H 'Authorization: Basic YWRtaW46YWJjMTIz' \
	-H 'Cache-Control: no-cache' \
	-H 'Connection: keep-alive' \
	-H 'Content-Type: text/plain;charset=UTF-8' \
	-H 'Origin: http://localhost:8080' \
	-H 'Pragma: no-cache' \
	-H 'Referer: http://localhost:8080/admin/config-public-details/' \
	-H 'Sec-Fetch-Dest: empty' \
	-H 'Sec-Fetch-Mode: cors' \
	-H 'Sec-Fetch-Site: same-origin' \
	-H 'Sec-GPC: 1' \
	-H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36' \
	--data-raw '{"value":"# Header\n\nLorem ipsum dolor sit amet, _consectetur adipiscing elit_, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Proin sagittis nisl rhoncus mattis rhoncus. **Mattis pellentesque id nibh tortor id aliquet lectus proin nibh**. Justo donec enim diam vulputate ut pharetra sit amet aliquam. Consectetur lorem donec massa sapien faucibus et molestie. Nisi scelerisque eu ultrices vitae auctor eu augue ut lectus. Lectus magna fringilla urna porttitor rhoncus. A cras semper auctor neque vitae. Sodales neque sodales ut etiam. Interdum velit euismod in pellentesque massa placerat duis ultricies lacus. Dignissim enim sit amet venenatis urna cursus eget.\n\n| Head | Head | Head |\n| --- | --- | --- |\n| Data | Data | Data |\n| Data | Data | Data |\n\n<img src=\"https://pixabay.com/get/g6693fc23171b84946e6533f87ca08dca39dfae8badf1ee1d659e15dbacaa893b18d1e378db050e8a6bcf156b741fb2be21b1c3834d16b2462549cea4e1621229_640.jpg\"/>\n\n## Subheader\n\nUllamcorper morbi tincidunt ornare massa eget egestas. Quis hendrerit dolor magna eget est. Sit amet tellus cras adipiscing enim eu turpis. Pharetra convallis posuere morbi leo urna molestie at. Commodo quis imperdiet massa tincidunt nunc pulvinar. Risus viverra adipiscing at in tellus integer. Felis donec et odio pellentesque diam volutpat. Enim nulla aliquet porttitor lacus luctus accumsan tortor posuere ac. Vel eros donec ac odio tempor orci. Turpis in eu mi bibendum neque egestas. Turpis egestas maecenas pharetra convallis posuere morbi. Non nisi est sit amet facilisis magna etiam tempor. Bibendum ut tristique et egestas quis ipsum suspendisse. Vel fringilla est ullamcorper eget nulla facilisi.\n\n1. List item 1\n2. List item 2\n3. List item 3\n\n[This is a link](https://owncast.online)\n\nScelerisque mauris pellentesque pulvinar pellentesque habitant morbi tristique. Sed vulputate mi sit amet mauris commodo quis imperdiet massa. Fermentum posuere urna nec tincidunt. Ut lectus arcu bibendum at varius vel pharetra vel. Ultricies mi quis hendrerit dolor magna eget. Rutrum tellus pellentesque eu tincidunt. Venenatis a condimentum vitae sapien pellentesque habitant. Eget nullam non nisi est sit amet facilisis. Tincidunt vitae semper quis lectus nulla at volutpat diam ut."}' \
	--compressed

# Links

curl 'http://localhost:8080/api/admin/config/socialhandles' \
	-H 'Accept: */*' \
	-H 'Accept-Language: en-US,en;q=0.9' \
	-H 'Authorization: Basic YWRtaW46YWJjMTIz' \
	-H 'Cache-Control: no-cache' \
	-H 'Connection: keep-alive' \
	-H 'Content-Type: text/plain;charset=UTF-8' \
	-H 'Origin: http://localhost:8080' \
	-H 'Pragma: no-cache' \
	-H 'Referer: http://localhost:8080/admin/config-public-details/' \
	-H 'Sec-Fetch-Dest: empty' \
	-H 'Sec-Fetch-Mode: cors' \
	-H 'Sec-Fetch-Site: same-origin' \
	-H 'Sec-GPC: 1' \
	-H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36' \
	--data-raw '{"value":[{"platform":"github","url":"https://github.com/owncast/owncast"},{"url":"https://facebook.biz/facebook","platform":"facebook"},{"url":"https://linkedin.biz/linkedin","platform":"linkedin"},{"url":"https://twitter.biz/twitter","platform":"twitter"}]}' \
	--compressed

# External action buttons

curl 'http://localhost:8080/api/admin/config/externalactions' \
	-H 'Accept: */*' \
	-H 'Accept-Language: en-US,en;q=0.9' \
	-H 'Authorization: Basic YWRtaW46YWJjMTIz' \
	-H 'Cache-Control: no-cache' \
	-H 'Connection: keep-alive' \
	-H 'Content-Type: text/plain;charset=UTF-8' \
	-H 'Origin: http://localhost:8080' \
	-H 'Pragma: no-cache' \
	-H 'Referer: http://localhost:8080/admin/actions/' \
	-H 'Sec-Fetch-Dest: empty' \
	-H 'Sec-Fetch-Mode: cors' \
	-H 'Sec-Fetch-Site: same-origin' \
	-H 'Sec-GPC: 1' \
	-H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36' \
	--data-raw '{"value":[{"url":"https://owncast.online/docs","title":"Documentation","description":"","icon":"","color":"","openExternally":false},{"url":"https://media1.giphy.com/media/Ju7l5y9osyymQ/giphy.gif?cid=ecf05e47otegqpl9mqfz880pi861dnjm5loy6kyrquy9lku0&rid=giphy.gif&ct=g","title":"Important","description":"","icon":"https://findicons.com/files/icons/1184/quickpix_2008/128/rick_roll_d.png","color":"#c87fd7","openExternally":false},{"url":"https://randommer.io/random-images","title":"New Tab","description":"","icon":"","color":"","openExternally":true}]}' \
	--compressed

# Offline message

curl 'http://localhost:8080/api/admin/config/offlinemessage' \
	-H 'Accept: */*' \
	-H 'Accept-Language: en-US,en;q=0.9' \
	-H 'Authorization: Basic YWRtaW46YWJjMTIz' \
	-H 'Cache-Control: no-cache' \
	-H 'Connection: keep-alive' \
	-H 'Content-Type: text/plain;charset=UTF-8' \
	-H 'Origin: http://localhost:8080' \
	-H 'Pragma: no-cache' \
	-H 'Referer: http://localhost:8080/admin/config-public-details/' \
	-H 'Sec-Fetch-Dest: empty' \
	-H 'Sec-Fetch-Mode: cors' \
	-H 'Sec-Fetch-Site: same-origin' \
	-H 'Sec-GPC: 1' \
	-H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36' \
	--data-raw '{"value":"Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. I am offline. This is my message."}' \
	--compressed


# Chat welcome message

curl 'http://localhost:8080/api/admin/config/welcomemessage' \
  -H 'accept: */*' \
  -H 'accept-language: en-US,en;q=0.9' \
  -H 'authorization: Basic YWRtaW46aG05dTl4a2c=' \
  -H 'cache-control: no-cache' \
  -H 'content-type: text/plain;charset=UTF-8' \
  -H 'Origin: http://localhost:8080' \
  -H 'pragma: no-cache' \
  -H 'referer: http://localhost:8080/admin/config-chat/' \
  -H 'sec-fetch-dest: empty' \
  -H 'sec-fetch-mode: cors' \
  -H 'sec-fetch-site: same-origin' \
  -H 'sec-gpc: 1' \
  -H 'user-agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36' \
  --data-raw '{"value":"This is an example chat welcome message."}' \
  --compressed 