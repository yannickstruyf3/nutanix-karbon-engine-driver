Kontainer Engine Karbon Driver
===============================

**Note**: This driver is still under development and only Karbon Development clusters can be requests as a result.

**Note**: This is NOT a Nutanix supported plugin.

**Note**: The plugin works for both Karbon 2.0 and Karbon 2.1.

This document will describe how to use the Rancher Kontainer Engine Karbon Driver. This custom driver has been written for Nutanix Karbon 2.0 and 2.1. 
This driver will add Nutanix Karbon as a hosted Kubernetes provider to your Rancher instance.

The UI plugin code can be found in the `ui-cluster-driver-nutanix` folder.

Special thanks to https://github.com/tuxtof for showing me the customization possibilities Rancher has to offer and brainstorming on the idea.

# Prerequisites
* Have working Rancher 2.0 installation
* Karbon 2.0/2.1 enabled on Prism Central
* Optional: use a webserver for hosting the UI files

For the steps on how to host the UI on your own webserver, please go to following link and refer to the `Building` section: 
https://github.com/yannickstruyf3/ui-cluster-driver-nutanix/blob/master/README.md

# Installation

Login on your rancher instance and go to `Tools` and select `Drivers`. 
Click `Add Cluster Driver` and fill in the details below.

![alt text](https://github.com/yannickstruyf3/kontainer-engine-driver-karbon/raw/master/images/0_clusterdriver.png " ")

You can use the pre-compiled binaries and files on following URLs:
Download URL:
https://github.com/yannickstruyf3/kontainer-engine-driver-karbon/releases/download/2.1.0/kontainer-engine-driver-karbon-linux-amd64

Custom UI URL:
https://github.com/yannickstruyf3/ui-cluster-driver-nutanix/releases/download/2.1.0/component.js

Whitelist Domains:
github.com

If you are hosting the files on your own webserver, you need to enter the URLs applicable for your environment.

![alt text](https://github.com/yannickstruyf3/kontainer-engine-driver-karbon/raw/master/images/1_createdriver.png " ")

After pressing `Create`, please wait until the state next to the `Karbon` cluster driver is `Active`. 
![alt text](https://github.com/yannickstruyf3/kontainer-engine-driver-karbon/raw/master/images/2_activedriver.png " ")

The driver is now installed and ready to use.

# Usage

## Requesting a Karbon cluster
Go to `Clusters` and press `Add Cluster`
![alt text](https://github.com/yannickstruyf3/kontainer-engine-driver-karbon/raw/master/images/3_emptyclusterscreen.png " ")

If the Karbon driver was correctly installed and activated, you will see an additional option in the `Hosted Kubernetes Provider` section. Clicking the `Karbon` icon will open the request page.

![alt text](https://github.com/yannickstruyf3/kontainer-engine-driver-karbon/raw/master/images/4_driverselection.png " ")

Enter all of the details applicable for your environment and press `Create`.

**Note**: Currently there is no input validation on the request page. This is work in progress

![alt text](https://github.com/yannickstruyf3/kontainer-engine-driver-karbon/raw/master/images/5_requestcluster.png " ")

The cluster is now provisioning. This can be seen from Rancher and from Prism Central.


![alt text](https://github.com/yannickstruyf3/kontainer-engine-driver-karbon/raw/master/images/6_clusterrancherprovisioning.png " ") ![alt text](https://github.com/yannickstruyf3/kontainer-engine-driver-karbon/raw/master/images/7_clusterprismprovisioning.png " ")

Time to grab a coffee and let the automation do its job!

When the cluster is provisioned (marked as `Active`) you can see the total amount of nodes (including masters) and also the CPU and memory usage. 

![alt text](https://github.com/yannickstruyf3/kontainer-engine-driver-karbon/raw/master/images/8_clusteractive.png " ")

Clicking on the cluster will give all of the required information regarding the Karbon Cluster.

![alt text](https://github.com/yannickstruyf3/kontainer-engine-driver-karbon/raw/master/images/9_clusterdashboard.png " ")

![alt text](https://github.com/yannickstruyf3/kontainer-engine-driver-karbon/raw/master/images/10_nodeoverview.png " ")

## Scaling

The driver also implements Karbon cluster scaling. Go to your cluster and press `edit`.

![alt text](https://github.com/yannickstruyf3/kontainer-engine-driver-karbon/raw/master/images/11_edit.png " ")

The request page will open. All of the fields are read-only except `Amount of Worker Nodes`. You can modify this value to your desired amount of Worker nodes in your Karbon cluster. Lowering the amount of nodes will scale down the cluster, increasing the amount of nodes will of course perform a scale up. Press `Create` when ready.

![alt text](https://github.com/yannickstruyf3/kontainer-engine-driver-karbon/raw/master/images/12_editamountofnodes.png " ")

Cluster is now in updating state. This can be seen from Rancher and Prism Central

![alt text](https://github.com/yannickstruyf3/kontainer-engine-driver-karbon/raw/master/images/13_clusterupdatingstaterancher.png " ")
![alt text](https://github.com/yannickstruyf3/kontainer-engine-driver-karbon/raw/master/images/14_clusterupdatingstateprism.png " ")

When the update process is complete, Rancher will report on the new amount of nodes.

![alt text](https://github.com/yannickstruyf3/kontainer-engine-driver-karbon/raw/master/images/15_addednodes.png " ")


## Destroying

Deleting of the cluster is also supported. Just press `Delete` and confirm the deletion process.

![alt text](https://github.com/yannickstruyf3/kontainer-engine-driver-karbon/raw/master/images/16_clusterdelete.png " ")

# DIY
## Building Karbon Cluster driver

`make`

Will output driver binaries into the `dist` directory, these can be imported 
directly into Rancher and used as cluster drivers.  They must be distributed 
via URLs that your Rancher instance can establish a connection to and download 
the driver binaries.  For example, this driver is distributed via a GitHub 
release and can be downloaded from one of those URLs directly.

## UI Development
The UI code is located in the `ui-cluster-driver-nutanix` folder.

This package contains a small web-server that will serve up the custom driver UI at `http://localhost:3000/component.js`.  You can run this while developing and point the Rancher settings there.
* `npm install`
* `npm start`
* The driver name can be optionally overridden: `npm start -- --name=DRIVERNAME`
* The compiled files are viewable at http://localhost:3000.

## UI Building

For other users to see your driver, you need to build it and host the output on a server accessible from their browsers.

* `npm run build`
* Copy the contents of the `dist` directory onto a webserver.
  * If your Rancher is configured to use HA or SSL, the server must also be available via HTTPS.
