# Copyright 2017 The Nuclio Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import datetime
import json

events_log_file_path = '/tmp/events.json'

def handler(context, event):
    """post event to the request recorder"""

    if event.trigger.kind != 'http' or event.get_header('x-nuclio-invoke-trigger'):
        body = event.body.decode('utf-8')
        context.logger.info('Received event body: {0}'.format(body))

        # serialized record
        serialized_record = json.dumps({
            'body': body,
            'headers': dict(event.headers),
            'timestamp': datetime.datetime.utcnow().isoformat(),
        })

        # store in log file
        with open(events_log_file_path, 'a') as events_log_file:
            events_log_file.write(serialized_record + ', ')

    else:

        # read the log file
        try:
            with open(events_log_file_path, 'r') as events_log_file:
                events_log_file_contents = events_log_file.read()
        except IOError:
            events_log_file_contents = ''

        # make this valid JSON by removing last two chars (, ) and enclosing in [ ]
        encoded_event_log = '[' + events_log_file_contents[:-2] + ']'

        context.logger.info('Returning events: {0}'.format(encoded_event_log))

        # return json.loads(encoded_event_log)
        return encoded_event_log
