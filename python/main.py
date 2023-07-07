import snowflake.connector
import time

from dotdict import DotDict
from structs import WarehouseStructs
from utils import WarehouseSize


def main():
    #  _____ ____  _  _____    _     _____ ____  _____
    # /  __//  _ \/ \/__ __\  / \ /|/  __//  __\/  __/
    # |  \  | | \|| |  / \    | |_|||  \  |  \/||  \
    # |  /_ | |_/|| |  | |    | | |||  /_ |    /|  /_
    # \____\\____/\_/  \_/    \_/ \|\____\\_/\_\\____\

    username: str = ""
    account: str = ""
    role: str = ""
    warehouse_usage: str = ""
    warehouse_autoscale: str = ""
    min_size: str = "large"
    max_size: str = "xxxlarge"
    # how many queued query to be considered as high traffic, recommended to keep at 5
    queued_threshold: str = 5
    # upsize at 15 and downsize at 0
    #   hence 5 at default gives a more budget approach
    #   on handling sizing event
    default_queue_checkpoint: int = 5

    #  ____  _____  ____  ____
    # / ___\/__ __\/  _ \/  __\
    # |    \  / \  | / \||  \/|
    # \___ |  | |  | \_/||  __/
    # \____/  \_/  \____/\_/

    queue_checkpoint: int = default_queue_checkpoint

    while True:
        with snowflake.connector.connect(
            user=username,
            account=account,
            role=role,
            warehouse=warehouse_usage,
            authenticator="externalbrowser",
            session_parameters={
                "query_tag": "snowflake-was",
            },
            autocommit=True,
        ) as conn:
            with conn.cursor() as cur:
                cur.execute(f"show warehouses like '{warehouse_autoscale}';")
                columns = [col[0] for col in cur.description]
                rows: list[str] = [
                    {k: v for k, v in zip(columns, row)} for row in cur.fetchall()
                ]
                if len(rows) > 1:
                    raise ValueError(
                        f"Only 1 warehouse is allowed. Found [{len(rows)}] warehouses."
                    )
                show_wh_result: WarehouseStructs = DotDict(rows[0])

                warehouse_size: WarehouseSize = WarehouseSize(
                    size=show_wh_result.size, min_size=min_size, max_size=max_size
                )

                # escalation point algo
                # if queue is high, don't bother to downsize
                if show_wh_result.queued >= queued_threshold:
                    queue_checkpoint = min(queue_checkpoint + 1, 10)
                    print(f"checkpoint hit {queue_checkpoint}")
                # if queue is not high, don't bother to upsize
                elif show_wh_result.queued < queued_threshold:
                    queue_checkpoint = max(queue_checkpoint - 1, 0)
                    print(f"checkpoint hit {queue_checkpoint}")

                sizing_event_triggered: bool = False
                sizing_event: str = ""
                if queue_checkpoint == 15:
                    sizing_event_triggered = warehouse_size.upsize()
                    sizing_event = "upsizing"
                    queue_checkpoint = default_queue_checkpoint
                elif queue_checkpoint == 0:
                    sizing_event_triggered = warehouse_size.downsize()
                    sizing_event = "downsizing"
                    queue_checkpoint = default_queue_checkpoint

                if sizing_event_triggered:
                    try:
                        print(f"{sizing_event} warehouse to [{warehouse_size.size}]")
                        cur.execute(
                            f"alter warehouse {warehouse_autoscale} set warehouse_size = {warehouse_size.size};"
                        )
                    except snowflake.connector.errors.ProgrammingError as e:
                        raise snowflake.connector.errors.ProgrammingError(
                            f"Error {e.errno} ({e.sqlstate}): {e.msg} ({e.sfqid})"
                        )
        time.sleep(60)


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        print("Stopped")
