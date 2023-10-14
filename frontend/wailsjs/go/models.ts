export namespace main {
	
	export class ProjectInfo {
	    projectName: string;
	    projectPath: string;
	    vscodeWorkspaceStoragePath: string;
	    workspaceRecommendationsIgnore: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ProjectInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.projectName = source["projectName"];
	        this.projectPath = source["projectPath"];
	        this.vscodeWorkspaceStoragePath = source["vscodeWorkspaceStoragePath"];
	        this.workspaceRecommendationsIgnore = source["workspaceRecommendationsIgnore"];
	    }
	}

}

