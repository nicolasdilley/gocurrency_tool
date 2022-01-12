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

## agg backend is used to create plot as a .pdf file
mpl.use('agg')

import matplotlib.pyplot as plt

 # create folders

if not os.path.exists("./normal"):
    os.makedirs("./normal")
if not os.path.exists("./median"):
    os.makedirs("./median")
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
if not os.path.exists("./median/absolute"):
    os.makedirs("./median/absolute")
if not os.path.exists("./cloc"):
    os.makedirs("./cloc")
if not os.path.exists("./csv/cloc"):
    os.makedirs("./csv/cloc")
if not os.path.exists("./csv/chan"):
    os.makedirs("./csv/chan")
if not os.path.exists("./csv/wg"):
    os.makedirs("./csv/wg")
if not os.path.exists("./csv/mutex"):
    os.makedirs("./csv/mutex")

# Normal relative (Counting features per kPLOC)
normal_chan_relative_cloc = np.loadtxt(open("../stats/results/normal/relative.csv","r+"),
                         usecols=(1,2,3,4,5,6),
                         unpack = True,
                         delimiter = ',',
                         dtype = float
                         )

normal_chan_relative_cloc_all_fig, ax1 = plt.subplots()
ax1.boxplot(x=[ normal_chan_relative_cloc[0],
                normal_chan_relative_cloc[1],
                normal_chan_relative_cloc[2],
                normal_chan_relative_cloc[3],
                normal_chan_relative_cloc[4],
                normal_chan_relative_cloc[5]],
                labels=["chan","send","receive","select","close","range"])
ax1.set_ylim(top=50)
plt.ylabel("Number of features / kPLOC")

normal_chan_relative_cloc_all_fig.savefig('./normal/relative/all_fig.pdf',dpi=900)
plt.close(normal_chan_relative_cloc_all_fig)


normal_chan_relative_cloc_df = df(data = {
'chan' : normal_chan_relative_cloc[0]
,'send' : normal_chan_relative_cloc[1]
,'receive': normal_chan_relative_cloc[2]
,'select': normal_chan_relative_cloc[3]
,'close': normal_chan_relative_cloc[4]
,'range': normal_chan_relative_cloc[5]
})
open("./csv/chan/relative_cloc.csv", "w+").write(normal_chan_relative_cloc_df.describe().to_csv())



# ------------------ ####

# All absolute 
normal_all_absolute = np.loadtxt(open("../stats/results/normal/all_absolute.csv","r+"),
                         usecols=(1,2,3,4,5,6,7,8,9,10,11,12,13),
                         unpack = True,
                         delimiter = ',',
                         dtype = float
                         )


normal_all_absolute_all_fig, ax1 = plt.subplots(figsize=(11,8))


ax1.boxplot(x=[ normal_all_absolute[0],
                normal_all_absolute[1],
                normal_all_absolute[2],
                normal_all_absolute[3],
                normal_all_absolute[4],
                normal_all_absolute[5],
                normal_all_absolute[6],
                normal_all_absolute[7],
                normal_all_absolute[8],
                normal_all_absolute[9],
                normal_all_absolute[10],
                normal_all_absolute[11],
                normal_all_absolute[12]],
                labels=["chan","send","receive","select","close","range",
                "Waitgroup","Add(x)","Done()","Wait()","Mutex","Lock()","Unlock()"],showmeans=True)
plt.ylabel("Number of features")
plt.vlines([6.5,10.5],0,1400)
plt.yscale('log')
plt.autoscale(True)


normal_all_absolute_all_fig.savefig('./normal/absolute/all_data_fig.pdf')
plt.close(normal_all_absolute_all_fig)

normal_all_absolute_df = df(data = {
'chan' : normal_all_absolute[0]
,'send' : normal_all_absolute[1]
,'receive': normal_all_absolute[2]
,'select': normal_all_absolute[3]
,'close': normal_all_absolute[4]
,'range': normal_all_absolute[5]
,'Waigroup': normal_all_absolute[6]
,'Add(x)': normal_all_absolute[7]
,'Done()': normal_all_absolute[8]
,'Wait()': normal_all_absolute[9]
,'Mutex' : normal_all_absolute[10]
,'Lock' : normal_all_absolute[11]
,'Unlock': normal_all_absolute[12]
})
open("./csv/all_absolute.csv", "w+").write(normal_all_absolute_df.describe().to_csv())



# Normal absolute (Meaning averaging features without counting only featured files)
normal_absolute = np.loadtxt(open("../stats/results/normal/absolute.csv","r+"),
                         usecols=(1,2,3,4,5,6),
                         unpack = True,
                         delimiter = ',',
                         dtype = float
                         )


normal_absolute_all_fig, ax1 = plt.subplots()
ax1.set_ylim(top=1450)

ax1.boxplot(x=[ normal_absolute[0],
                normal_absolute[1],
                normal_absolute[2],
                normal_absolute[3],
                normal_absolute[4],
                normal_absolute[5]],
                labels=["chan","send","receive","select","close","range"])
plt.ylabel("Number of features")


normal_absolute_all_fig.savefig('./normal/absolute/all_fig.pdf',dpi=900)
plt.close(normal_absolute_all_fig)

normal_new_chan_df = df(data = {
'chan' : normal_absolute[0]
,'send' : normal_absolute[1]
,'receive': normal_absolute[2]
,'select': normal_absolute[3]
,'close': normal_absolute[4]
,'range': normal_absolute[5]
})
open("./csv/chan/absolute.csv", "w+").write(normal_new_chan_df.describe().to_csv())



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
ax1.set_ylim(top=8)
normal_new_chan_all_fig.savefig('./normal/new_chan/all_fig.pdf',dpi=900)
plt.close(normal_new_chan_all_fig)

# Save the 5 number as csv file also

normal_new_chan_df = df(data = {
'send' : normal_new_chan[0]
,'receive': normal_new_chan[1]
,'select': normal_new_chan[2]
,'close': normal_new_chan[3]
,'range': normal_new_chan[4]
})
open("./csv/chan/relative_chan.csv", "w+").write(normal_new_chan_df.describe().to_csv())


#  new_wg (counter relative to wg)
normal_new_wg = np.loadtxt(open("../stats/results/normal/new_wg.csv","r+"),
                         usecols=(1,2,3),
                         unpack = True,
                         delimiter = ',',
                         dtype = float
                         )

normal_new_wg_all_fig, ax1 = plt.subplots()
ax1.boxplot(x=[ normal_new_wg[0],
                normal_new_wg[1],
                normal_new_wg[2],
                ],
                labels=["Add(x)","Done()","Wait()"])
plt.ylabel("Number of features / Waitgroup")
ax1.set_ylim(top=10,bottom=-2)
normal_new_wg_all_fig.savefig('./normal/new_wg/all_fig.pdf',dpi=900)
plt.close(normal_new_wg_all_fig)


normal_new_wg_df = df(data = {
'Add(x)': normal_new_wg[0]
,'Done()': normal_new_wg[1]
,'Wait()': normal_new_wg[2]
})
open("./csv/wg/relative_mu.csv", "w+").write(normal_new_wg_df.describe().to_csv())


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
ax1.set_ylim(top=12.5,bottom=-2)
normal_new_mu_all_fig.savefig('./normal/new_mu/all_fig.pdf',dpi=900)
plt.close(normal_new_mu_all_fig)

normal_new_mutex_df = df(data = {
'Lock' : normal_new_mu[0]
,'Unlock': normal_new_mu[1]
})
open("./csv/mutex/relative_mu.csv", "w+").write(normal_new_mutex_df.describe().to_csv())

# All relative

normal_all_relative_all_fig, ax1 = plt.subplots(figsize=(11,8))

ax1.boxplot(x=[ normal_new_chan[0],
                normal_new_chan[1],
                normal_new_chan[2],
                normal_new_chan[3],
                normal_new_chan[4],
                normal_new_wg[0],
                normal_new_wg[1],
                normal_new_wg[2],
                normal_new_mu[0],
                normal_new_mu[1],],
                labels=["send","receive","select","close","range",
                "Add(x)","Done()","Wait()","Lock()","Unlock()"],)
plt.ylabel("Number of features / concurrency mechanisms")
plt.vlines([5.5,8.5],0,9)


normal_all_relative_all_fig.savefig('./normal/relative/relative_all_data_fig.pdf')
plt.close(normal_all_relative_all_fig)



# ------------------ ####

# Normal relative count for WG ! (with respect to kLoc)
normal_wg_relative_cloc = np.loadtxt(open("../stats/results/normal/relative_wg.csv","r+"),
                         usecols=(1,2,3,4),
                         unpack = True,
                         delimiter = ',',
                         dtype = float
                         )

normal_wg_relative_cloc_all_fig, ax1 = plt.subplots()
ax1.boxplot(x=[ normal_wg_relative_cloc[0],
                normal_wg_relative_cloc[1],
                normal_wg_relative_cloc[2],
                normal_wg_relative_cloc[3],
                ],
                labels=["Waitgroup","Add(x)","Done()","Wait()"])
ax1.set_ylim(top=50)
plt.ylabel("Number of features / kPLOC")

normal_wg_relative_cloc_all_fig.savefig('./normal/relative/all_fig_wg.pdf',dpi=900)
plt.close(normal_wg_relative_cloc_all_fig)

normal_wg_relative_cloc_df = df(data = {
'Waitgroup' : normal_wg_relative_cloc[0]
,'Add(x)': normal_wg_relative_cloc[1]
,'Done()': normal_wg_relative_cloc[2]
,'Wait()': normal_wg_relative_cloc[3]
})
open("./csv/wg/relative_cloc.csv", "w+").write(normal_wg_relative_cloc_df.describe().to_csv())



# ------------------ ####

# Normal relative count for Mutex ! (with respect to kLoc)
normal_mu_relative_cloc = np.loadtxt(open("../stats/results/normal/relative_mu.csv","r+"),
                         usecols=(1,2,3),
                         unpack = True,
                         delimiter = ',',
                         dtype = float
                         )

normal_mu_relative_cloc_all_fig, ax1 = plt.subplots()
ax1.boxplot(x=[ normal_mu_relative_cloc[0],
                normal_mu_relative_cloc[1],
                normal_mu_relative_cloc[2],
                ],
                labels=["Mutex","Lock","Unlock"])
ax1.set_ylim(top=20)
plt.ylabel("Number of features / kPLOC")

normal_mu_relative_cloc_all_fig.savefig('./normal/relative/all_fig_mu.pdf',dpi=900)
plt.close(normal_mu_relative_cloc_all_fig)

normal_mu_relative_cloc_df = df(data = {
'Mutex' : normal_mu_relative_cloc[0]
,'Lock' : normal_mu_relative_cloc[1]
,'Unlock': normal_mu_relative_cloc[2]
})

open("./csv/mutex/relative_cloc.csv", "w+").write(normal_mu_relative_cloc_df.describe().to_csv())

# ------------------ ####

# All relative to cloc ! 

normal_all_relative_cloc, ax1 = plt.subplots(figsize=(11,8))


ax1.boxplot(x=[ normal_chan_relative_cloc[0],
                normal_chan_relative_cloc[1],
                normal_chan_relative_cloc[2],
                normal_chan_relative_cloc[3],
                normal_chan_relative_cloc[4],
                normal_chan_relative_cloc[5],
                normal_wg_relative_cloc[0],
                normal_wg_relative_cloc[1],
                normal_wg_relative_cloc[2],
                normal_wg_relative_cloc[3],
                normal_mu_relative_cloc[0],
                normal_mu_relative_cloc[1],
                normal_mu_relative_cloc[2]],
                labels=["chan","send","receive","select","close","range",
                "Waitgroup","Add(x)","Done()","Wait()","Mutex","Lock()","Unlock()"],showmeans=True)
plt.ylabel("Number of features / kPLOC")
plt.vlines([6.5,10.5],0,35)
ax1.set_yscale('log')

normal_all_relative_cloc.savefig('./normal/relative/all_cloc.pdf')
plt.close(normal_all_relative_cloc)

normal_absolute_all_fig, ax1 = plt.subplots()
ax1.set_ylim(top=250)
ax1.boxplot(x=[ normal_absolute[0],
                normal_absolute[1],
                normal_absolute[2],
                normal_absolute[3],
                normal_absolute[4],],
                labels=["Waitgroup","Add(const)","Add(x)","Done()","Wait()"])
plt.ylabel("Number of features")
normal_absolute_all_fig.savefig('./normal/absolute/all_fig_wg.pdf',dpi=900)
plt.close(normal_absolute_all_fig)

normal_new_wg_df = df(data = {
'Waitgroup' : normal_absolute[0]
,'Add(const)' : normal_absolute[1]
,'Add(x)': normal_absolute[2]
,'Done()': normal_absolute[3]
,'Wait()': normal_absolute[4]
})
open("./csv/wg/absolute.csv", "w+").write(normal_new_wg_df.describe().to_csv())

# =============== *** ==============


# Normal absolute count for Mutex ! absolute number of features in projects
normal_absolute = np.loadtxt(open("../stats/results/normal/absolute_mu.csv","r+"),
                         usecols=(1,2,3),
                         unpack = True,
                         delimiter = ',',
                         dtype = float
                         )

normal_absolute_all_fig, ax1 = plt.subplots()
ax1.set_ylim(top=120)

ax1.boxplot(x=[ normal_absolute[0],
                normal_absolute[1],
                normal_absolute[2],
                ],
                labels=["Mutex","Lock","Unlock"])
plt.ylabel("Number of features")
normal_absolute_all_fig.savefig('./normal/absolute/all_fig_mu.pdf',dpi=900)
plt.close(normal_absolute_all_fig)


normal_new_wg_df = df(data = {
'Mutex' : normal_absolute[0]
,'Lock' : normal_absolute[1]
,'Unlock': normal_absolute[2]
})
open("./csv/mutex/absolute.csv", "w+").write(normal_new_wg_df.describe().to_csv())


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
normal_absolute_cloc_loc.savefig('./cloc/normal_absolute_cloc_loc.pdf',dpi=900)
plt.close(normal_absolute_cloc_loc)

normal_absolute_all_fig, ax1 = plt.subplots(squeeze=True)
ax1.boxplot(x=[normal_absolute[0],
                normal_absolute[1],
                normal_absolute[2]],
                labels=["size","package","file"])
plt.ylabel("Percent")
normal_absolute_all_fig.savefig('./cloc/normal_absolute_all_fig.pdf',dpi=900)
plt.close(normal_absolute_all_fig)

normal_absolute_table = df(data = {
'size' : normal_absolute[0]
,'packages': normal_absolute[1]
,'files': normal_absolute[2]
})
open("./csv/cloc/cloc_normal_absolute.csv", "w").write(normal_absolute_table.describe().to_csv())


# Median CLOC ffiles informations


median_absolute = np.loadtxt(open("../stats/results/cloc/median.csv","r+"),
                         usecols=(1,2,3),
                         unpack = True,
                         delimiter = ',',
                         dtype = float
                         )
median_absolute_cloc_loc, ax1 = plt.subplots(figsize=(4,9))
ax1.boxplot(x=median_absolute[0],
                labels=["CLOC"],widths=[0.4])
# plt.subplots_adjust(left=0.7, right=0.9)
plt.ylabel("Percent")
plt.tight_layout()
median_absolute_cloc_loc.savefig('./cloc/median_absolute_cloc_loc.pdf',dpi=900)
plt.close(median_absolute_cloc_loc)

median_absolute_all_fig, ax1 = plt.subplots(squeeze=True)
ax1.boxplot(x=[median_absolute[0],
                median_absolute[1],
                median_absolute[2]],
                labels=["size","package","file"])
plt.ylabel("Percent")
median_absolute_all_fig.savefig('./cloc/median_absolute_all_fig.pdf',dpi=900)
plt.close(median_absolute_all_fig)

median_absolute_table = df(data = {
'size' : median_absolute[0]
,'packages': median_absolute[1]
,'files': median_absolute[2]
})
open("./csv/cloc/cloc_median_absolute.csv", "w").write(median_absolute_table.describe().to_csv())



### MEDIAN 
median_all_absolute = np.loadtxt(open("../stats/results/median/all_absolute.csv","r+"),
                         usecols=(1,2,3,4,5,6,7,8,9,10,11,12,13),
                         unpack = True,
                         delimiter = ',',
                         dtype = float
                         )


median_all_absolute_all_fig, ax1 = plt.subplots(figsize=(11,8))
# ax1.set_ylim(top=1450)
ax1.boxplot(x=[ median_all_absolute[0],
                median_all_absolute[1],
                median_all_absolute[2],
                median_all_absolute[3],
                median_all_absolute[4],
                median_all_absolute[5],
                median_all_absolute[6],
                median_all_absolute[7],
                median_all_absolute[8],
                median_all_absolute[9],
                median_all_absolute[10],
                median_all_absolute[11],
                median_all_absolute[12]],
                labels=["chan","send","receive","select","close","range",
                "Waitgroup","Add(x)","Done()","Wait()","Mutex","Lock()","Unlock()"],showmeans=True)
plt.ylabel("Number of features")
plt.yscale('log')
plt.vlines([6.5,10.5],0,60)


median_all_absolute_all_fig.savefig('./median/absolute/all_data_fig.pdf')
plt.close(median_all_absolute_all_fig)

median_all_absolute_df = df(data = {
'chan' : median_all_absolute[0]
,'send' : median_all_absolute[1]
,'receive': median_all_absolute[2]
,'select': median_all_absolute[3]
,'close': median_all_absolute[4]
,'range': median_all_absolute[5]
,'Waigroup': median_all_absolute[6]
,'Add(x)': median_all_absolute[7]
,'Done()': median_all_absolute[8]
,'Wait()': median_all_absolute[9]
,'Mutex' : median_all_absolute[10]
,'Lock' : median_all_absolute[11]
,'Unlock': median_all_absolute[12]
})
open("./csv/median_all_absolute.csv", "w+").write(median_all_absolute_df.describe().to_csv())



# median absolute (Meaning averaging features without counting only featured files)
median_absolute = np.loadtxt(open("../stats/results/median/absolute.csv","r+"),
                         usecols=(1,2,3,4,5,6),
                         unpack = True,
                         delimiter = ',',
                         dtype = float
                         )


median_absolute_all_fig, ax1 = plt.subplots()
ax1.boxplot(x=[ median_absolute[0],
                median_absolute[1],
                median_absolute[2],
                median_absolute[3],
                median_absolute[4],
                median_absolute[5]],
                labels=["chan","send","receive","select","close","range"])
plt.ylabel("Number of features")


median_absolute_all_fig.savefig('./median/absolute/all_fig.pdf',dpi=900)
plt.close(median_absolute_all_fig)

median_new_chan_df = df(data = {
'chan' : median_absolute[0]
,'send' : median_absolute[1]
,'receive': median_absolute[2]
,'select': median_absolute[3]
,'close': median_absolute[4]
,'range': median_absolute[5]
})

open("./csv/chan/median_absolute.csv", "w+").write(median_new_chan_df.describe().to_csv())



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
known_size_chan_fig.savefig('./known_size_chan_fig.pdf',dpi=900)
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
non_zero_known_chan_fig.savefig('./non_zero_known_chan_fig.pdf',dpi=900)
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
go_in_any_for_fig.savefig('./go_in_any_for_fig.pdf',dpi=900)
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
chan_in_any_for_fig.savefig('./chan_in_any_for_fig.pdf',dpi=900)
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
go_in_for_fig.savefig('./go_in_for_fig.pdf',dpi=900)
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
chan_in_for_fig.savefig('./chan_in_for_fig.pdf',dpi=900)
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
go_in_constant_for_fig.savefig('./go_in_constant_for_fig.pdf',dpi=900)
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
go_per_projects_fig.savefig('./go_per_projects_fig.pdf',dpi=900)
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
constant_go_in_constant_for_fig.savefig('./constant_go_in_constant_for_fig.pdf',dpi=900)
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
num_branch_per_projects_fig.savefig('./num_branch_per_projects_fig.pdf',dpi=900)
plt.close(num_branch_per_projects_fig)

num_branch_per_projects_table = df(data = {
'Goroutines' : num_branch_per_projects
})

open("./csv/num_branch_per_projects.csv", "w").write(num_branch_per_projects_table.describe().to_csv())