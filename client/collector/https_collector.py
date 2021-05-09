import asyncio
from pyobe import Probe, ProbeOptions

options = ProbeOptions()
options.name = "test"
options.centre_url = "http://localhost:12000"
options.concurrency = 64
options.timeout = 30


class BaseProbe(Probe):
    def __init__(self, path="https://"):
        super(BaseProbe, self).__init__(options)
        self.path = path

    async def gen_url(self):
        url = f"{self.centre_url}/task?path={self.path}"
        while True:
            yield url

    @staticmethod
    def get_target_url(base_url):
        return base_url

    async def on_url(self, url: str):
        while True:
            res = await self.session.get(url)
            j = await res.json()
            task_url = self.get_target_url(j.get("url"))
            tid = j.get("id")
            if url is None or tid is None:
                await asyncio.sleep(10)
            else:
                break
        res = await self.session.get(task_url,
                                     proxy="http://192.168.1.119:8118")
        if res.status >= 300:
            await asyncio.sleep(1)
            return
        data = await res.read()
        raw = data.decode("utf-8-sig")

        await self.session.post(
            f"{self.centre_url}/raw",
            json={"data": raw, "task_id": tid, "url": task_url},
        )


if __name__ == "__main__":
    BaseProbe().start()
