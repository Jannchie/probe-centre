import logging
import requests

BASE_URL = "https://api.bilibili.com/x/web-interface/ranking/v2?rid=%d&type=%s"


def add_rank():
    rank_types = ["rookie", "origin", "all"]
    tid_list = [0, 5, 181, 155, 119, 217, 211, 160, 188, 36, 4, 129, 3, 1, 168]
    for rank_type in rank_types:
        for tid in tid_list:
            url = BASE_URL % (tid, rank_type)
            res = requests.post(
                "http://localhost:12000/task",
                json={"url": url},
            )
            print(res.json())


def parse_rank():
    url = "http://localhost:12000/raw?" \
          "type=json&" \
          "count=200&" \
          "path=https://api.bilibili.com/x/web-interface/ranking/v2&consume=1"
    res = requests.get(url)
    data = res.json()
    for json_data in data:
        d = json_data['data']
        data_list = d.get('data').get('list')
        for video_data in data_list:
            try:
                aid = video_data.get('aid')
                mid = video_data.get('owner').get('mid')
                requests.post("http://localhost:12000/task", json=dict(
                    url=f"https://api.bilibili.com/x/space/acc/info?mid={mid}"))
                requests.post("http://localhost:12000/task", json=dict(
                    url=f"https://api.bilibili.com/x/relation/stat?vmid={mid}"))
                requests.post("http://localhost:12000/task", json=dict(
                    url=f"https://api.bilibili.com/x/ugcpay-rank/elec"
                        f"/month/up?up_mid={mid}"))
                requests.post("http://localhost:12000/task", json=dict(
                    url=f"https://api.bilibili.com/x/web-interface/view?aid={aid}"))
            except Exception as e:
                logging.exception(e)


if __name__ == "__main__":
    # add_rank()
    parse_rank()
