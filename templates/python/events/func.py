from cloudevents.http import CloudEvent, to_binary

# ---------------------------------------
# Function template
# context is a dictionary with the Flask request object and 
# and a cloud_event key for incoming events
# ---------------------------------------
def main(context):
    # print(f"Method: {context['request'].method}")
    attributes = {
      "type": "com.example.fn",
      "source": "https://example.com/fn"
    }
    data = { "message": "Howdy!" }
    event = CloudEvent(attributes, data)
    headers, body = to_binary(event)
    return body, 200, headers