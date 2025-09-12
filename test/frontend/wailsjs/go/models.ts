export namespace domain {
	
	export class Task {
	    id: number;
	    task: string;
	    priority: string;
	    status: boolean;
	    due_date: string;
	
	    static createFrom(source: any = {}) {
	        return new Task(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.task = source["task"];
	        this.priority = source["priority"];
	        this.status = source["status"];
	        this.due_date = source["due_date"];
	    }
	}

}

