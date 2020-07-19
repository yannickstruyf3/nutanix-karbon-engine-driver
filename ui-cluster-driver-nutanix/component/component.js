/*!!!!!!!!!!!Do not change anything between here (the DRIVERNAME placeholder will be automatically replaced at buildtime)!!!!!!!!!!!*/
// https://github.com/rancher/ui/blob/master/lib/shared/addon/mixins/cluster-driver.js
import ClusterDriver from 'shared/mixins/cluster-driver';

// do not remove LAYOUT, it is replaced at build time with a base64 representation of the template of the hbs template
// we do this to avoid converting template to a js file that returns a string and the cors issues that would come along with that
const LAYOUT;
/*!!!!!!!!!!!DO NOT CHANGE END!!!!!!!!!!!*/


/*!!!!!!!!!!!GLOBAL CONST START!!!!!!!!!!!*/
// EMBER API Access - if you need access to any of the Ember API's add them here in the same manner rather then import them via modules, since the dependencies exist in rancher we dont want to expor the modules in the amd def
const computed = Ember.computed;
const observer = Ember.observer;
const get = Ember.get;
const set = Ember.set;
const alias = Ember.computed.alias;
const service = Ember.inject.service;
const all = Ember.RSVP.all;

/*!!!!!!!!!!!GLOBAL CONST END!!!!!!!!!!!*/

const reclaimPolicyMap = {
  'Delete': 'Delete',
  'Retain': 'Retain',
}

const cniProviderMap = {
  'Flannel': 'Flannel',
  'Calico': 'Calico',
}
const filesystemMap = {
  'ext4': 'ext4',
  'xfs': 'xfs',
}

const karbonVersionAndChoicesMap = {
  '2.1': {
    '1.16.10-0': '1.16.10-0',
    '1.15.12-0': '1.15.12-0',
    '1.14.10-1': '1.14.10-1',
    '1.13.12-1': '1.13.12-1',
  },
  '2.0': {
    '1.16.8-0': '1.16.8-0',
    '1.15.10-0': '1.15.10-0',
    '1.14.10-0': '1.14.10-0',
    '1.13.12-0': '1.13.12-0',
  },
}


const deploymentOptionsMap = {
  'Development': {
    'amountOfMasters': {
      "1": 1,
    },
    'amountOfETCDs': {
      "1": 1,
    }
  },
  'Production - active/passive': {
    'amountOfMasters': {
      "2": 2,
    },
    'amountOfETCDs': {
      "3": 3,
      "5": 5,
    }
  },
  'Production - active/active': {
    'amountOfMasters': {
      "2": 2,
      "3": 3,
      "4": 4,
      "5": 5,
    },
    'amountOfETCDs': {
      "3": 3,
      "5": 5,
    }
  },
}



/*!!!!!!!!!!!DO NOT CHANGE START!!!!!!!!!!!*/
export default Ember.Component.extend(ClusterDriver, {
  driverName: '%%DRIVERNAME%%',
  configField: '%%DRIVERNAME%%EngineConfig', // 'googleKubernetesEngineConfig'
  app: service(),
  router: service(),
  /*!!!!!!!!!!!DO NOT CHANGE END!!!!!!!!!!!*/

  visibleCheckMasterIp3: function () {
    return get(this, 'cluster.%%DRIVERNAME%%EngineConfig.masternodes') > 2;
  }.property('cluster.%%DRIVERNAME%%EngineConfig.masternodes'),
  visibleCheckMasterIp4: function () {
    return get(this, 'cluster.%%DRIVERNAME%%EngineConfig.masternodes') > 3;
  }.property('cluster.%%DRIVERNAME%%EngineConfig.masternodes'),
  visibleCheckMasterIp5: function () {
    return get(this, 'cluster.%%DRIVERNAME%%EngineConfig.masternodes') > 4;
  }.property('cluster.%%DRIVERNAME%%EngineConfig.masternodes'),




  init() {
    /*!!!!!!!!!!!DO NOT CHANGE START!!!!!!!!!!!*/
    // This does on the fly template compiling, if you mess with this :cry:
    const decodedLayout = window.atob(LAYOUT);
    const template = Ember.HTMLBars.compile(decodedLayout, {
      moduleName: 'shared/components/cluster-driver/driver-%%DRIVERNAME%%/template'
    });
    set(this, 'layout', template);

    this._super(...arguments);
    /*!!!!!!!!!!!DO NOT CHANGE END!!!!!!!!!!!*/

    let config = get(this, 'config');
    let configField = get(this, 'configField');


    if (!config) {
      config = this.get('globalStore').createRecord({
        type: configField,
        host: "",
        port: 9440,
        username: "",
        password: "",
        insecure: false,
        workernodes: 1,
        image: "",
        cniprovider: "Flannel",
        version: "1.16.10-0",
        cluster: "",
        vmnetwork: "",
        workercpu: 8,
        workermemorymib: 8192,
        workerdiskmib: 122880,
        mastercpu: 2,
        mastermemorymib: 4096,
        masterdiskmib: 122880,
        etcdcpu: 4,
        etcdmemorymib: 8192,
        etcddiskmib: 40960,
        clusteruser: "",
        clusterpassword: "",
        storagecontainer: "",
        filesystem: "ext4",
        reclaimpolicy: "Delete",
        flashmode: false,
        deployment: "Development",
        masternodes: 1,
        etcdnodes: 1,
        mastervipip: "",
        masterip1: "",
        masterip2: "",
        masterip3: "",
        masterip4: "",
        masterip5: "",
      });

      set(this, 'cluster.%%DRIVERNAME%%EngineConfig', config);
    }
  },

  config: alias('cluster.%%DRIVERNAME%%EngineConfig'),


  actions: {
    save(cb) {
      if (!this.validate()) {
        cb(false)
      }
      this.send('driverSave');
    },
    cancel() {
      console.log("Cancel")
      // probably should not remove this as its what every other driver uses to get back
      get(this, 'router').transitionTo('global-admin.clusters.index');
    },
  },

  // Add custom validation beyond what can be done from the config API schema


  validate() {
    // Get generic API validation errors
    this._super();
    var errors = get(this, 'errors') || [];
    if (!get(this, 'cluster.name')) {
      errors.push('Name is required');
    }

    // Check if values are set
    let keys = {
      "host": { name: "Prism Central Host", type: "string" },
      "port": { name: "Prism Central Port", type: "integer" },
      "username": { name: "Prism Central Username", type: "string" },
      "password": { name: "Prism Central Password", type: "string" },
      "workernodes": { name: "Amount of Worker Nodes", type: "integer" },
      "image": { name: "Kubernetes Image", type: "string" },
      "version": { name: "Kubernetes Version", type: "string" },
      "cluster": { name: "Nutanix Cluster", type: "string" },
      "vmnetwork": { name: "VM Network", type: "string" },
      "workercpu": { name: "Worker CPU", type: "integer" },
      "workermemorymib": { name: "Worker Memory", type: "integer" },
      "workerdiskmib": { name: "Worker Storage", type: "integer" },
      "mastercpu": { name: "Master CPU", type: "integer" },
      "mastermemorymib": { name: "Master Memory", type: "integer" },
      "masterdiskmib": { name: "Master Storage", type: "integer" },
      "etcdcpu": { name: "ETCD CPU", type: "integer" },
      "etcdmemorymib": { name: "ETCD Memory", type: "integer" },
      "etcddiskmib": { name: "ETCD Storage", type: "integer" },
      "clusteruser": { name: "Prism Element User", type: "string" },
      "clusterpassword": { name: "Prism Element Password", type: "string" },
      "storagecontainer": { name: "Storage Container", type: "string" },
      "filesystem": { name: "Filesystem", type: "string" },
      "reclaimpolicy": { name: "Reclaim Policy", type: "string" },
      "deployment": { name: "Deployment Type", type: "string" },
      "masternodes": { name: "Amount of Master Nodes", type: "integer" },
      "etcdnodes": { name: "Amount of ETCD Nodes", type: "integer" },
    }

    for (var k in keys) {
      let val = keys[k]
      let config = get(this, 'cluster.%%DRIVERNAME%%EngineConfig.' + k)
      if (!config) {
        errors.push(val['name'] + ' is required');
      }
      if (val['type'] == 'integer') {
        if (config < 1) {
          errors.push('Value of ' + val['name'] + ' must be >= 1');
        }
      }
    }

    //check if dropdown values are respected
    if (!["ext4", "xfs"].includes(get(this, 'cluster.%%DRIVERNAME%%EngineConfig.filesystem'))) {
      errors.push('Filesystem must be ext4 or xfs');
    }
    if (!["Delete", "Retain"].includes(get(this, 'cluster.%%DRIVERNAME%%EngineConfig.reclaimpolicy'))) {
      errors.push('Reclaim policy must be Delete or Retain');
    }
    if (!["Deployment", 'Production - active/passive', 'Production - active/active'].includes(get(this, 'cluster.%%DRIVERNAME%%EngineConfig.deployment'))) {
      errors.push('Deployment type must be "Deployment", "Production - active/passive" or "Production - active/active"');
    }
    // Add more specific errors
    // Set the array of errors for display,
    // and return true if saving should continue.
    if (get(errors, 'length')) {
      set(this, 'errors', errors);
      return false;
    } else {
      set(this, 'errors', null);
      return true;
    }
  },

  // Any computed properties or custom logic can go here
  reclaimPolicyChoices: Object.entries(reclaimPolicyMap).map((e) => ({
    label: e[1],
    value: e[0]
  })),

  cniProviderChoices: Object.entries(cniProviderMap).map((e) => ({
    label: e[1],
    value: e[0]
  })),

  filesystemChoices: Object.entries(filesystemMap).map((e) => ({
    label: e[1],
    value: e[0]
  })),
  versionChoices: computed('cluster.%%DRIVERNAME%%EngineConfig.karbonversion', function () {
    return Object.entries(karbonVersionAndChoicesMap["2.1"]).map((e) => ({
      label: e[0],
      value: e[0]
    }))
  }),

  deploymentChoices: Object.entries(deploymentOptionsMap).map((e) => ({
    label: e[0],
    value: e[0]
  })),
  amountOfMasterNodesChoices: computed('cluster.%%DRIVERNAME%%EngineConfig.deployment', function () {
    let deployment = get(this, 'cluster.%%DRIVERNAME%%EngineConfig.deployment');
    return Object.entries(deploymentOptionsMap[deployment]["amountOfMasters"]).map((e) => ({
      label: e[0],
      value: e[1]
    }))
  }),
  amountOfETCDNodesChoices: computed('cluster.%%DRIVERNAME%%EngineConfig.deployment', function () {
    let deployment = get(this, 'cluster.%%DRIVERNAME%%EngineConfig.deployment');
    return Object.entries(deploymentOptionsMap[deployment]["amountOfETCDs"]).map((e) => ({
      label: e[0],
      value: e[1]
    }))
  }),

});

