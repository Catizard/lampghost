export namespace entity {
	
	export class RivalInfo {
	
	
	    static createFrom(source: any = {}) {
	        return new RivalInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}

}

