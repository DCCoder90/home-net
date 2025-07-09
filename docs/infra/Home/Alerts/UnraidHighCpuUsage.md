# UnraidHighCpuUsage

```yaml
apiVersion: 1
groups:
    - orgId: 1
      name: 1MEval
      folder: UnraidAlerts
      interval: 1m
      rules:
        - uid: denqatodeumm8f
          title: UnraidHighCpuUsage
          condition: C
          data:
            - refId: A
              relativeTimeRange:
                from: 600
                to: 0
              datasourceUid: aengwcm978mpsb
              model:
                disableTextWrap: false
                editorMode: code
                expr: 100-(avg(rate(node_cpu_seconds_total{mode="idle"}[5m]))*100)
                fullMetaSearch: false
                includeNullMetadata: true
                instant: true
                intervalMs: 1000
                legendFormat: __auto
                maxDataPoints: 43200
                range: false
                refId: A
                useBackend: false
            - refId: C
              datasourceUid: __expr__
              model:
                conditions:
                    - evaluator:
                        params:
                            - 90
                        type: gt
                      operator:
                        type: and
                      query:
                        params:
                            - C
                      reducer:
                        params: []
                        type: last
                      type: query
                datasource:
                    type: __expr__
                    uid: __expr__
                expression: A
                intervalMs: 1000
                maxDataPoints: 43200
                refId: C
                type: threshold
          noDataState: NoData
          execErrState: Error
          for: 5m
          annotations:
            description: CPU usage is {{ $value | printf "%.2f" }}% (averaged over 5m) for the last 5 minutes.
            summary: High CPU usage on Unraid server
          labels: {}
          isPaused: false
          notification_settings:
            receiver: Discord Hook
        - uid: eenqc2ayyae4gb
          title: Disk2SpaceWarning
          condition: C
          data:
            - refId: A
              relativeTimeRange:
                from: 600
                to: 0
              datasourceUid: aengwcm978mpsb
              model:
                disableTextWrap: false
                editorMode: code
                expr: (node_filesystem_free_bytes{mountpoint="/mnt/disk2"}/node_filesystem_size_bytes{mountpoint="/mnt/disk2"})*100
                fullMetaSearch: false
                includeNullMetadata: true
                instant: true
                intervalMs: 1000
                legendFormat: __auto
                maxDataPoints: 43200
                range: false
                refId: A
                useBackend: false
            - refId: C
              datasourceUid: __expr__
              model:
                conditions:
                    - evaluator:
                        params:
                            - 95
                        type: gt
                      operator:
                        type: and
                      query:
                        params:
                            - C
                      reducer:
                        params: []
                        type: last
                      type: query
                datasource:
                    type: __expr__
                    uid: __expr__
                expression: A
                intervalMs: 1000
                maxDataPoints: 43200
                refId: C
                type: threshold
          noDataState: NoData
          execErrState: Error
          for: 1m
          annotations:
            description: Unraid disk 2 usage is at {{$value|printf"%.2f"}}%.
            summary: Disk 2 space critically low!
          isPaused: false
          notification_settings:
            receiver: Discord Hook
        - uid: cenqcf89he0owa
          title: Disk1SpaceWarning
          condition: C
          data:
            - refId: A
              relativeTimeRange:
                from: 600
                to: 0
              datasourceUid: aengwcm978mpsb
              model:
                disableTextWrap: false
                editorMode: code
                expr: (node_filesystem_free_bytes{mountpoint="/mnt/disk1"}/node_filesystem_size_bytes{mountpoint="/mnt/disk1"})*100
                fullMetaSearch: false
                includeNullMetadata: true
                instant: true
                intervalMs: 1000
                legendFormat: __auto
                maxDataPoints: 43200
                range: false
                refId: A
                useBackend: false
            - refId: C
              datasourceUid: __expr__
              model:
                conditions:
                    - evaluator:
                        params:
                            - 95
                        type: gt
                      operator:
                        type: and
                      query:
                        params:
                            - C
                      reducer:
                        params: []
                        type: last
                      type: query
                datasource:
                    type: __expr__
                    uid: __expr__
                expression: A
                intervalMs: 1000
                maxDataPoints: 43200
                refId: C
                type: threshold
          noDataState: NoData
          execErrState: Error
          for: 1m
          annotations:
            description: Unraid disk 1 usage is at {{$value|printf"%.2f"}}%.
            summary: Disk 1 space critically low!
          isPaused: false
          notification_settings:
            receiver: Discord Hook
        - uid: denqchiib17nke
          title: Disk3SpaceWarning
          condition: C
          data:
            - refId: A
              relativeTimeRange:
                from: 600
                to: 0
              datasourceUid: aengwcm978mpsb
              model:
                disableTextWrap: false
                editorMode: code
                expr: (node_filesystem_free_bytes{mountpoint="/mnt/disk3"}/node_filesystem_size_bytes{mountpoint="/mnt/disk3"})*100
                fullMetaSearch: false
                includeNullMetadata: true
                instant: true
                intervalMs: 1000
                legendFormat: __auto
                maxDataPoints: 43200
                range: false
                refId: A
                useBackend: false
            - refId: C
              datasourceUid: __expr__
              model:
                conditions:
                    - evaluator:
                        params:
                            - 95
                        type: gt
                      operator:
                        type: and
                      query:
                        params:
                            - C
                      reducer:
                        params: []
                        type: last
                      type: query
                datasource:
                    type: __expr__
                    uid: __expr__
                expression: A
                intervalMs: 1000
                maxDataPoints: 43200
                refId: C
                type: threshold
          noDataState: NoData
          execErrState: Error
          for: 1m
          annotations:
            description: Unraid disk 3 usage is at {{$value|printf"%.2f"}}%.
            summary: Disk 3 space critically low!
          isPaused: false
          notification_settings:
            receiver: Discord Hook
```