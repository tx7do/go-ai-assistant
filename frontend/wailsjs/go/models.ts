export namespace azopenai {
	
	export class AssistantConfig {
	    assistantId?: string;
	    fileId?: string;
	    threadId?: string;
	    name?: string;
	    description?: string;
	    instructions?: string;
	
	    static createFrom(source: any = {}) {
	        return new AssistantConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.assistantId = source["assistantId"];
	        this.fileId = source["fileId"];
	        this.threadId = source["threadId"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.instructions = source["instructions"];
	    }
	}

}

