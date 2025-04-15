import matplotlib.pyplot as plt
import itertools
import os
import geopandas as gpd
from shapely.geometry import Point
import contextily as cx
from tqdm import tqdm  

# Get the directory where this script lives
script_dir = os.path.dirname(os.path.realpath(__file__))

# Create patterns_output folder inside the script's directory
output_dir = os.path.join(script_dir, "patterns_output")
os.makedirs(output_dir, exist_ok=True)

# File path for the text file (in the same script folder)
combo_file_path = os.path.join(script_dir, "patterns_nodes.txt")

# Node Name : (lat, lon)
pos = {
    "T6": (28.52836141567019, 77.57777379378554),
    "R": (28.527471246410833, 77.57890336254033),
    "A": (28.52692140098942, 77.57706407672785),
    "CnD": (28.525519393291614, 77.57653461322104),
    "G": (28.52799850829152, 77.5749038301257),
    "Arc": (28.52723498480195, 77.57292972432113),
    "SARC": (28.523582249548888, 77.57437275274097),
    "C1": (28.52440905859402, 77.57308311017349),
    "Dib": (28.525235547619324, 77.57072373132858),
    "DH3": (28.52316003957018, 77.5696713403205),
    "ISC": (28.521496152454514, 77.5712575598366),
}

nodes = list(pos.keys())
c3 = list(itertools.combinations(nodes, 3))
c4 = list(itertools.combinations(nodes, 4))
total = c3 + c4

# Write all combinations to patterns_nodes.txt
with open(combo_file_path, "w") as f:
    f.write(str(total))

# Convert all nodes to GeoDataFrame and reproject to Web Mercator
pts = gpd.GeoDataFrame(
    {"name": nodes},
    geometry = [Point(lon, lat) for lat, lon in pos.values()],
    crs = "EPSG:4326"
).to_crs(epsg=3857)

# Plot and save each pattern
for i, combo in enumerate(tqdm(total, desc="Generating patterns")):
    fig, ax = plt.subplots(figsize=(8, 8))
    ax.set_title(f"Nodes: {', '.join(combo)}", fontsize=8)

#     # Set map limits around campus
    buffer = 200
    bounds = pts.total_bounds
    ax.set_xlim(bounds[0] - buffer, bounds[2] + buffer)
    ax.set_ylim(bounds[1] - buffer, bounds[3] + buffer)

#     # Add basemap (with transparency)
    cx.add_basemap(ax, source=cx.providers.OpenStreetMap.Mapnik, alpha=0.5)

#     # Plot only the selected nodes as red triangles
    for name in combo:
        geom = pts[pts["name"] == name].geometry.values[0]
        ax.scatter(geom.x, geom.y, s=100, color='red', edgecolors='black', marker='^', zorder=3)
        ax.text(
            geom.x, geom.y - 25, name,
            ha='center', va='top',
            fontsize=10,
            zorder=4
        )

    ax.axis("off")
    plt.tight_layout()
    filename = f"pattern_{i:03d}"  # Use numbered filenames
    plt.savefig(os.path.join(output_dir, f"{filename}.png"), dpi=300)
    plt.close()

print("done")
