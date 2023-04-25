Query to create a table for booking:

create table project.booking (
bookingId int,	
seatId	int,	
name	varchar(255),
contact	int,
PRIMARY KEY (bookingId),
FOREIGN KEY (seatId) REFERENCES seatList(id)
);
