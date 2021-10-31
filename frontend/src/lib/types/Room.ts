export interface PublicRoom {
	id: string;
	options: string[];
	question: string;
}

export interface Room extends PublicRoom {
	ownerID: string;
}
