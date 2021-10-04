import numpy as np
from pandas import DataFrame as df
from pandas.plotting import table
from scipy.stats import trim_mean, kurtosis
from scipy.stats.mstats import mode, gmean, hmean
import matplotlib as mpl
import csv
import string
import os
from matplotlib.ticker import ScalarFormatter

## agg backend is used to create plot as a .png file
mpl.use('agg')

import matplotlib.pyplot as plt

 # create folders

if not os.path.exists("./normal"):
    os.makedirs("./normal")
if not os.path.exists("./normal/relative"):
    os.makedirs("./normal/relative")
if not os.path.exists("./normal/absolute"):
    os.makedirs("./normal/absolute")
if not os.path.exists("./normal/new_chan"):
    os.makedirs("./normal/new_chan")
if not os.path.exists("./normal/new_wg"):
    os.makedirs("./normal/new_wg")
if not os.path.exists("./normal/new_mu"):
    os.makedirs("./normal/new_mu")
if not os.path.exists("./cloc"):
    os.makedirs("./cloc")
if not os.path.exists("./csv/cloc"):
    os.makedirs("./csv/cloc")

# Normal relative (Counting features per kPLOC)
normal_relative = np.loadtxt(open("../stats/results/normal/relative.csv","r+"),
                         usecols=(1,2,3,4,5,6),
                         unpack = True,
                         delimiter = ',',
                         dtype = float
                         )

normal_relative_all_fig, ax1 = plt.subplots()
ax1.boxplot(x=[ normal_relative[0],
                normal_relative[1],
                normal_relative[2],
                normal_relative[3],
                normal_relative[4],
                normal_relative[5]],
                labels=["chan","send","receive","select","close","range"])
ax1.set_ylim(top=83)
plt.ylabel("Number of features / kPLOC")

normal_relative_all_fig.savefig('./normal/relative/all_fig.png',dpi=900)
plt.close(normal_relative_all_fig)


# ------------------ ####


# Normal absolute (Meaning averaging features without counting only featured files)
normal_absolute = np.loadtxt(open("../stats/results/normal/absolute.csv","r+"),
                         usecols=(1,2,3,4,5,6),
                         unpack = True,
                         delimiter = ',',
                         dtype = float
                         )


normal_absolute_all_fig, ax1 = plt.subplots()
ax1.boxplot(x=[ normal_absolute[0],
                normal_absolute[1],
                normal_absolute[2],
                normal_absolute[3],
                normal_absolute[4],
                normal_absolute[5]],
                labels=["chan","send","receive","select","close","range"])
plt.ylabel("Number of features")
normal_absolute_all_fig.savefig('./normal/absolute/all_fig.png',dpi=900)
plt.close(normal_absolute_all_fig)


# ------------------ ####

# Normal new_chan (counter relative to chan)
normal_new_chan = np.loadtxt(open("../stats/results/normal/new_chan.csv","r+"),
                         usecols=(1,2,3,4,5),
                         unpack = True,
                         delimiter = ',',
                         dtype = float
                         )

normal_new_chan_all_fig, ax1 = plt.subplots()
ax1.boxplot(x=[ normal_new_chan[0],
                normal_new_chan[1],
                normal_new_chan[2],
                normal_new_chan[3],
                normal_new_chan[4]],
                labels=["send","receive","select","close","range"])
plt.ylabel("Number of features / channels")
ax1.set_ylim(top=20,bottom=-2)
normal_new_chan_all_fig.savefig('./normal/new_chan/all_fig.png',dpi=900)
plt.close(normal_new_chan_all_fig)

#  new_wg (counter relative to wg)
normal_new_wg = np.loadtxt(open("../stats/results/normal/new_wg.csv","r+"),
                         usecols=(1,2,3,4),
                         unpack = True,
                         delimiter = ',',
                         dtype = float
                         )

normal_new_wg_all_fig, ax1 = plt.subplots()
ax1.boxplot(x=[ normal_new_wg[0],
                normal_new_wg[1],
                normal_new_wg[2],
                normal_new_wg[3],
                ],
                labels=["Add(const)","Add(x)","Done()","Wait()"])
plt.ylabel("Number of features / Waitgroup")
ax1.set_ylim(top=20,bottom=-2)
normal_new_wg_all_fig.savefig('./normal/new_wg/all_fig.png',dpi=900)
plt.close(normal_new_wg_all_fig)

#  new_mu (counter relative to mu)
normal_new_mu = np.loadtxt(open("../stats/results/normal/new_mu.csv","r+"),
                         usecols=(1,2),
                         unpack = True,
                         delimiter = ',',
                         dtype = float
                         )

normal_new_mu_all_fig, ax1 = plt.subplots()
ax1.boxplot(x=[ normal_new_mu[0],
                normal_new_mu[1],
                ],
                labels=["Lock","Unlock"])
plt.ylabel("Number of features / Mutex")
ax1.set_ylim(top=20,bottom=-2)
normal_new_mu_all_fig.savefig('./normal/new_mu/all_fig.png',dpi=900)
plt.close(normal_new_mu_all_fig)

# ------------------ ####

# Normal relative count for WG ! (with respect to kLoc)
normal_relative = np.loadtxt(open("../stats/results/normal/relative_wg.csv","r+"),
                         usecols=(1,2,3,4,5),
                         unpack = True,
                         delimiter = ',',
                         dtype = float
                         )

normal_relative_all_fig, ax1 = plt.subplots()
ax1.boxplot(x=[ normal_relative[0],
                normal_relative[1],
                normal_relative[2],
                normal_relative[3],
                normal_relative[4],
                ],
                labels=["Waitgroup","Add(const)","Add(x)","Done()","Wait()"])
ax1.set_ylim(top=83)
plt.ylabel("Number of features / kPLOC")

normal_relative_all_fig.savefig('./normal/relative/all_fig_wg.png',dpi=900)
plt.close(normal_relative_all_fig)

# ------------------ ####


# Normal absolute count for WG ! absolute number of features in projects
normal_absolute = np.loadtxt(open("../stats/results/normal/absolute_wg.csv","r+"),
                         usecols=(1,2,3,4,5),
                         unpack = True,
                         delimiter = ',',
                         dtype = float
                         )

normal_absolute_all_fig, ax1 = plt.subplots()
ax1.boxplot(x=[ normal_absolute[0],
                normal_absolute[1],
                normal_absolute[2],
                normal_absolute[3],
                normal_absolute[4],],
                labels=["Waitgroup","Add(const)","Add(x)","Done()","Wait()"])
plt.ylabel("Number of features")
normal_absolute_all_fig.savefig('./normal/absolute/all_fig_wg.png',dpi=900)
plt.close(normal_absolute_all_fig)


# =============== *** ==============


# Normal absolute count for Mutex ! absolute number of features in projects
normal_absolute = np.loadtxt(open("../stats/results/normal/absolute_mu.csv","r+"),
                         usecols=(1,2,3),
                         unpack = True,
                         delimiter = ',',
                         dtype = float
                         )

normal_absolute_all_fig, ax1 = plt.subplots()
ax1.boxplot(x=[ normal_absolute[0],
                normal_absolute[1],
                normal_absolute[2],
                ],
                labels=["Mutex","Lock","Unlock"])
plt.ylabel("Number of features")
normal_absolute_all_fig.savefig('./normal/absolute/all_fig_mu.png',dpi=900)
plt.close(normal_absolute_all_fig)


# =============== *** ==============

# Normal relative count for Mutex ! relative number of features in projects / KLOC
normal_absolute = np.loadtxt(open("../stats/results/normal/relative_mu.csv","r+"),
                         usecols=(1,2,3),
                         unpack = True,
                         delimiter = ',',
                         dtype = float
                         )

normal_absolute_all_fig, ax1 = plt.subplots()
ax1.boxplot(x=[ normal_absolute[0],
                normal_absolute[1],
                normal_absolute[2],
                ],
                labels=["Mutex","Lock","Unlock"])
plt.ylabel("Number of features")
normal_absolute_all_fig.savefig('./normal/absolute/all_fig_mu.png',dpi=900)
plt.close(normal_absolute_all_fig)


# =============== *** ==============

normal_absolute = np.loadtxt(open("../stats/results/cloc/normal.csv","r+"),
                         usecols=(1,2,3),
                         unpack = True,
                         delimiter = ',',
                         dtype = float
                         )
normal_absolute_cloc_loc, ax1 = plt.subplots(figsize=(4,9))
ax1.boxplot(x=normal_absolute[0],
                labels=["CLOC"],widths=[0.4])
# plt.subplots_adjust(left=0.7, right=0.9)
plt.ylabel("Percent")
plt.tight_layout()
normal_absolute_cloc_loc.savefig('./cloc/normal_absolute_cloc_loc.png',dpi=900)
plt.close(normal_absolute_cloc_loc)

normal_absolute_all_fig, ax1 = plt.subplots(squeeze=True)
ax1.boxplot(x=[normal_absolute[0],
                normal_absolute[1],
                normal_absolute[2]],
                labels=["size","package","file"])
plt.ylabel("Percent")
normal_absolute_all_fig.savefig('./cloc/normal_absolute_all_fig.png',dpi=900)
plt.close(normal_absolute_all_fig)

normal_absolute_table = df(data = {
'size' : normal_absolute[0]
,'packages': normal_absolute[1]
,'files': normal_absolute[2]
})
open("./csv/cloc/cloc_normal_absolute.csv", "w").write(normal_absolute_table.describe().to_csv())

# CHAN SIZE TABLE


known_size_chan = np.loadtxt(open("../stats/results/known_size_chan.csv","r+"),
                         usecols=(1),
                         unpack = True,
                         delimiter = ',',
                         dtype = int
                         )
known_size_chan_fig, ax1 = plt.subplots(figsize=(4.5,9))
ax1.boxplot(x=known_size_chan,
                labels=["Size of channels"],widths=[0.4])
plt.ylabel("Size")
ax1.set_yscale('log')
known_size_chan_fig.savefig('./known_size_chan_fig.png',dpi=900)
plt.close(known_size_chan_fig)

known_size_chan_table = df(data = {
'known size' : known_size_chan
})
open("./csv/known_size_chan.csv", "w").write(known_size_chan_table.describe().to_csv())


non_zero_known_chan = np.loadtxt(open("../stats/results/non_zero_known_chan.csv","r+"),
                         usecols=(1),
                         unpack = True,
                         delimiter = ',',
                         dtype = int
                         )
non_zero_known_chan_fig, ax1 = plt.subplots(figsize=(4.7,9))
ax1.boxplot(non_zero_known_chan,
                labels=["Size of channels"],widths=[0.4])
ax1.set_yscale('log')

plt.ylabel("Size")
non_zero_known_chan_fig.savefig('./non_zero_known_chan_fig.png',dpi=900)
plt.close(non_zero_known_chan_fig)

non_zero_known_chan_table = df(data = {
'known size' : non_zero_known_chan
})
open("./csv/non_zero_known_chan.csv", "w").write(non_zero_known_chan_table.describe().to_csv())


go_in_any_for = np.loadtxt(open("../stats/results/go_in_any_for.csv","r+"),
                         usecols=(1),
                         unpack = True,
                         delimiter = ',',
                         dtype = float
                         )
go_in_any_for_fig, ax1 = plt.subplots(figsize=(4.7,9))
ax1.boxplot(go_in_any_for,
                labels=["Goroutines in any for"],widths=[0.4])
plt.ylabel("Number of goroutines")
go_in_any_for_fig.savefig('./go_in_any_for_fig.png',dpi=900)
plt.close(go_in_any_for_fig)
ax1.set_ylim(top=250)

go_in_any_for_table = df(data = {
'known size' : go_in_any_for
})
open("./csv/go_in_any_for.csv", "w").write(go_in_any_for_table.describe().to_csv())


chan_in_any_for = np.loadtxt(open("../stats/results/chan_in_any_for.csv","r+"),
                         usecols=(1),
                         unpack = True,
                         delimiter = ',',
                         dtype = float
                         )
chan_in_any_for_fig, ax1 = plt.subplots(figsize=(4.7,9))
ax1.boxplot(chan_in_any_for,
                labels=["Channels in any for"],widths=[0.4])
plt.ylabel("Number of chans")
chan_in_any_for_fig.savefig('./chan_in_any_for_fig.png',dpi=900)
plt.close(chan_in_any_for_fig)

chan_in_any_for_table = df(data = {
'known size' : chan_in_any_for
})
open("./csv/chan_in_any_for.csv", "w").write(chan_in_any_for_table.describe().to_csv())


go_in_for = np.loadtxt(open("../stats/results/go_in_unknown_for.csv","r+"),
                         usecols=(1),
                         unpack = True,
                         delimiter = ',',
                         dtype = float
                         )
go_in_for_fig, ax1 = plt.subplots(figsize=(4.7,9))
ax1.boxplot(go_in_for,
                labels=["Goroutines in For"],widths=[0.4])
plt.ylabel("Number of goroutines")
go_in_for_fig.savefig('./go_in_for_fig.png',dpi=900)
plt.close(go_in_for_fig)

go_in_for_table = df(data = {
'known size' : go_in_for
})
open("./csv/go_in_unknown_for.csv", "w").write(go_in_for_table.describe().to_csv())


chan_in_for = np.loadtxt(open("../stats/results/chan_in_for.csv","r+"),
                         usecols=(1),
                         unpack = True,
                         delimiter = ',',
                         dtype = float
                         )
chan_in_for_fig, ax1 = plt.subplots(figsize=(4.7,9))
ax1.boxplot(chan_in_for,
                labels=["Channels in for"],widths=[0.4])
plt.ylabel("Number of chans")
chan_in_for_fig.savefig('./chan_in_for_fig.png',dpi=900)
plt.close(chan_in_for_fig)

chan_in_for_table = df(data = {
'known size' : chan_in_for
})
open("./csv/chan_in_for.csv", "w").write(chan_in_for_table.describe().to_csv())

go_in_constant_for = np.loadtxt(open("../stats/results/go_in_constant_for.csv","r+"),
                         usecols=(1),
                         unpack = True,
                         delimiter = ',',
                         dtype = float
                         )
go_in_constant_for_fig, ax1 = plt.subplots(figsize=(4.7,9))
ax1.boxplot(go_in_constant_for,
                labels=["Goroutines in bounded for"],widths=[0.4])
plt.ylabel("Number of goroutines")
go_in_constant_for_fig.savefig('./go_in_constant_for_fig.png',dpi=900)
plt.close(go_in_constant_for_fig)

go_in_constant_for_table = df(data = {
'known size' : go_in_constant_for
})
open("./csv/go_in_constant_for.csv", "w").write(go_in_constant_for_table.describe().to_csv())

go_per_projects = np.loadtxt(open("../stats/results/go_per_projects.csv","r+"),
                         usecols=(1),
                         unpack = True,
                         delimiter = ',',
                         dtype = float
                         )
go_per_projects_fig, ax1 = plt.subplots(figsize=(4.7,9))
ax1.boxplot(go_per_projects,
                labels=["Goroutines"],widths=[0.4])
plt.ylabel("Number of goroutines")
go_per_projects_fig.savefig('./go_per_projects_fig.png',dpi=900)
plt.close(go_per_projects_fig)

go_per_projects_table = df(data = {
'Goroutines' : go_per_projects
})
open("./csv/go_per_projects.csv", "w").write(go_per_projects_table.describe().to_csv())

constant_go_in_constant_for = np.loadtxt(open("../stats/results/constant_go_in_constant_for.csv","r+"),
                         usecols=(1),
                         unpack = True,
                         delimiter = ',',
                         dtype = float
                         )
constant_go_in_constant_for_fig, ax1 = plt.subplots(figsize=(4.7,9))
ax1.boxplot(constant_go_in_constant_for,
                labels=["Size"],widths=[0.4])
plt.ylabel("Size of constant for")
constant_go_in_constant_for_fig.savefig('./constant_go_in_constant_for_fig.png',dpi=900)
plt.close(constant_go_in_constant_for_fig)

constant_go_in_constant_for_table = df(data = {
'Goroutines' : constant_go_in_constant_for
})
open("./csv/constant_go_in_constant_for.csv", "w").write(constant_go_in_constant_for_table.describe().to_csv())

num_branch_per_projects = np.loadtxt(open("../stats/results/num_branch_per_projects.csv","r+"),
                         usecols=(1),
                         unpack = True,
                         delimiter = ',',
                         dtype = float
                         )
num_branch_per_projects_fig, ax1 = plt.subplots(figsize=(4.7,9))
ax1.boxplot(num_branch_per_projects,
                labels=["Size"],widths=[0.4])
plt.ylabel("Size of constant for")
num_branch_per_projects_fig.savefig('./num_branch_per_projects_fig.png',dpi=900)
plt.close(num_branch_per_projects_fig)

num_branch_per_projects_table = df(data = {
'Goroutines' : num_branch_per_projects
})

open("./csv/num_branch_per_projects.csv", "w").write(num_branch_per_projects_table.describe().to_csv())