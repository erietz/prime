import subprocess
from pathlib import Path
import json
from matplotlib import pyplot as plt

low = 0
high = 100000
dest_dir = Path("./benchmark")


def run_benchmark():
    if not dest_dir.exists():
        dest_dir.mkdir()

    n_partitions = 1
    while n_partitions < high:
        inner_cmd = f'./prime --lo {low} --hi {high} --nPartitions {n_partitions}'

        cmd = [
            'hyperfine',
            '--export-json',
            f'./benchmark/lo_{low}_hi_{high}_nPartitions_{n_partitions}.json',
            inner_cmd
        ]

        print("running", cmd)
        subprocess.run(cmd)
        n_partitions *= 2


def load_data():
    data = []
    for file in dest_dir.glob("*.json"):
        with open(file) as f:
            data.append(json.load(f)["results"][0])
    return data


def plot():
    data = load_data()

    commands = [d["command"] for d in data]
    goroutines = [int(c.split(" ")[-1]) for c in commands]
    means = [d["mean"] for d in data]
    stddev = [d["stddev"] for d in data]
    # print(sorted(zip(*[goroutines, commands, means, stddev])))

    data = sorted(zip(*[goroutines, commands, means, stddev]))
    goroutines = [i[0] for i in data]
    commands = [i[1] for i in data]
    means = [i[2] for i in data]
    stddev = [i[3] for i in data]

    plt.bar(goroutines, means, yerr=stddev, width=200)
    plt.xlabel("goroutines")
    plt.ylabel("time / s")
    plt.savefig('./media/benchmark.png')
    plt.close()
    # plt.show()

    plt.bar(goroutines, means, yerr=stddev, width=1)
    plt.xlabel("goroutines")
    plt.ylabel("time / s")
    plt.xlim([0, 100])
    plt.savefig('./media/benchmark_zoom.png')
    plt.close()
    # plt.show()

    # plt.errorbar(goroutines, means, yerr=stddev, fmt='o')
    # plt.xlabel("goroutines")
    # plt.ylabel("time / s")
    # plt.xlim([0, 100])
    # plt.savefig('./media/benchmark_zoom_errorbar.png')
    # plt.close()
    # # plt.show()


# run_benchmark()
plot()
