export interface PublicRoom {
	id: string;
	options: PublicOption[];
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

export interface PublicOption {
	id: string;
	value: string;
	available: boolean;
	selectedByMe: boolean;
}

export interface Option extends PublicOption {
	id: string;
	value: string;
	available: boolean;
	selectedByMe: boolean;
	selectedByName?: string;
}
