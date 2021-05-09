from https_collector import BaseProbe


class BiliProbe(BaseProbe):
    def __init__(self):
        super().__init__("probe://bilibili")

    @staticmethod
    def get_target_url(base_url):
        return base_url


if __name__ == '__main__':
    BiliProbe().start()
