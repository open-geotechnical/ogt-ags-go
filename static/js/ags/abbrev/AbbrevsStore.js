
Ext.define('ags.abbrev.AbbrevsStore', {
	extend: 'Ext.data.Store',
	requires: [
       //'Ags.model.Abbrev'
    ],
	constructor: function(){
		Ext.apply(this, {
			model: 'ags.model.Abbrev',
			storeId: "abbrevs",
			sssorters: [ {
				property: 'dated',
				direction: 'DESC'
			}],
			deadsortInfo: {
				property: 'code',
				direction: 'DESC'
			},
			groupField: "group",
			pageSize: 2000,
			autoLoad: true,
			proxy: {
				type: 'ajax',
				url: "/ajax/ags4/abbreviations.json",
				reader: {
					type: 'json',
					root: "abbreviations",
					idProperty: 'head_code',
					totalProperty: 'abbreviations_count'
				}
			}
		});
		this.callParent();

	}
});
