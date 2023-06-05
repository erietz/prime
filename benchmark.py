import subprocess
from pathlib import Path
# import json
# from matplotlib import pyplot as plt

low = 0
high = 100000
dest_dir = Path("./benchmark")


def run_benchmark():
    if not dest_dir.exists():
        dest_dir.mkdir()

    n_partitions = 1
    while n_partitions < high:
        n_partitions *= 2

        inner_cmd = [
            './prime',
            '--lo',
            str(low),
            '--hi',
            str(high),
            '--nPartitions',
            str(n_partitions)
        ]

        cmd = [
            'hyperfine',
            '--export-json',
            f'./benchmark/lo_{low}_hi_{high}_nPartitions_{n_partitions}.json',
            ' '.join(inner_cmd)
        ]

        print("running", cmd)
        subprocess.run(cmd)
