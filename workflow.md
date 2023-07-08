# High Level
```mermaid
stateDiagram-v2
    state queued_checkpoint <<choice>>
    state queued_query_count <<choice>>
    curr: Current WH Size
    reset: reset back to default_queue_checkpoint
    query_gte: query >= equal queued_threshold
    query_lt: query < queued_threshold

    
    [*] --> curr
    curr --> queued_query_count
    queued_query_count --> query_lt
    queued_query_count --> query_gte
    query_lt --> queued_checkpoint: queue_checkpoint - 1
    query_gte --> queued_checkpoint: queue_checkpoint + 1
    queued_checkpoint --> downsize: current wh > min_size \n&& \nqueued_checkpoint < queued_threshold
    queued_checkpoint --> upsize: current wh < max_size \n&& \nqueued_checkpoint >= queued_threshold
    downsize --> reset
    upsize --> reset
```

# Example
Let's assume queue threshold is at 4, meaning when the warehouse autoscaler detects there are 4 queries queued, queued_checkpoint will increase by 1, vice versa. The following are the assumed configuration:

```yaml
username: jdoe@email.com
account: hello.us-east-1
role: wh_admin
warehouse_usage: compute_wh_admin
warehouse_autoscale: compute_wh_analysts
min_size: xsmall
max_size: xxlarge
current_size: xsmall
queued_threshold: 4
default_queue_checkpoint: 5
```

```mermaid
timeline
    title upsize event
        queued_checkpoint=5 : queued queries=5 : queued_checkpoint increases by 1
        ... :  assume queued queries is greater than or equal 4 (queued_threshold) until queued_checkpoint=14
        queued_checkpoint=15: queued queries=8: warehouse will upsize by one given the warehouse is below max_size
            : compute_wh_analysts = small now 
        queued_checkpoint=5 : queud_checkpoint is reset back to 5 after resizing
```

So after upsizing, the warehouse manage to start clearing the queued queries and now the queue has been reduced.

```mermaid
timeline
    title downsize event
        queued_checkpoint=5 : queued queries=3 : queued_checkpoint decreases by 1
        ... : assume queued queries is less than 4 (queued_threshold) until queued_checkpoint=1
        queued_checkpoint=0: queued queries=2: warehouse will downsize by one given the warehouse is above min_size
            : compute_wh_analysts = xsmall now 
        queued_checkpoint=5 : queud_checkpoint is reset back to 5 after resizing
```


# Disclaimer
As of current design, the upsize interval is hardcoded at 10 mins and downsize interval is hardcoded at 5 mins.