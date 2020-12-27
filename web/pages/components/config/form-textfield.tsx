/*
- auto saves ,ajax call
- set default text
- show error state/confirm states
- show info
- label
- min/max length

- populate with curren val (from local sstate)

load page, 
get all config vals, 
save to local state/context.
read vals from there.
update vals to state, andthru api.


*/
import React, { useContext } from 'react';
import { ServerStatusContext } from '../../../utils/server-status-context';



Server Name
<Input placeholder="Owncast" value={name} />