export interface PublicRoom {
	id: string;
	options: Option[];
	question: string;
	ownedByMe: boolean;
}

export interface Room {
	id: string;
	options: Option[];
	question: string;
}

export interface SimpleRoom {
	id: string;
	question: string;
}

export interface Option {
	id: string;
	value: string;
	available: boolean;
	selectedByMe: boolean;
}
