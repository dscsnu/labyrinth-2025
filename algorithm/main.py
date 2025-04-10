import matplotlib.pyplot as plt
import itertools
import os

# Get the directory where this script lives
script_dir = os.path.dirname(os.path.realpath(__file__))

# Create patterns_output folder inside the script's directory
output_dir = os.path.join(script_dir, "patterns_output")
os.makedirs(output_dir, exist_ok=True)

# File path for the text file (in the same script folder)
combo_file_path = os.path.join(script_dir, "patterns_nodes.txt")

# Node No : (x,y)
positions = {
    0: (0, 2),
    1: (1, 2),
    2: (2, 2),
    3: (0, 1),
    4: (1, 1),
    5: (2, 1),
    6: (0, 0),
    7: (1, 0),
    8: (2, 0),
    9: (1, -1),
}

# Generate combinations
combos_3 = list(itertools.combinations(range(10), 3))
combos_4 = list(itertools.combinations(range(10), 4))
all_combos = combos_3 + combos_4

# Write all combos to the patterns_nodes.txt
with open(combo_file_path, "w") as f:
    f.write(str(all_combos))

# Plot and save each one
for i, combo in enumerate(all_combos):
    fig, ax = plt.subplots(figsize=(4, 4))
    ax.set_title(f"Nodes: {combo}")
    ax.set_xlim(-1, 3)
    ax.set_ylim(-2, 3)
    ax.set_aspect('equal')
    ax.axis('off')

    for node, (x, y) in positions.items():
        color = 'gray' if node in combo else 'lightblue'
        size = 300 if node in combo else 200
        ax.scatter(x, y, s=size, color=color, edgecolors='black', zorder=3)
        ax.text(x, y, str(node), ha='center', va='center', zorder=4, fontsize=9)

    plt.tight_layout()
    plt.savefig(os.path.join(output_dir, f"pattern_{i:03d}.png"))
    plt.close()

print("done")
