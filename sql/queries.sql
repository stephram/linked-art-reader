select * from db_entities, db_identifiers where entity_id = "https://mart.getty.edu/museum/collection/person/da88f653-216e-4dd3-9816-eeeff36c81f3";
select * from db_entities, db_identifiers where uuid = "da88f653-216e-4dd3-9816-eeeff36c81f3";

select * from db_entities where uuid = "da88f653-216e-4dd3-9816-eeeff36c81f3";
select db_entities.id from db_entities where uuid = "da88f653-216e-4dd3-9816-eeeff36c81f3";
select * from db_identifiers where entity_ref = 1;
select identifier_id, label, value from db_entities, db_identifiers where uuid = "da88f653-216e-4dd3-9816-eeeff36c81f3";

select db_entities.id as entity_id, db_identifiers.label, db_identifiers.value from db_entities, db_identifiers inner join db_entities on entity_id = 1;

select count() from db_entities;
select count() from db_entities where type not in ('Person', 'Group');
select uuid, title from db_entities where type = 'HumanMadeObject';
select * from db_entities where type = 'HumanMadeObject' and title = 'Irises';
select * from db_classifiers where entity_ref = '51076';
select * from db_references where entity_ref = '51076';
select * from db_locations where entity_ref = '51076';

drop view entityview;
drop view entityview; 
create view entityview as;
select db_identifiers.entity_ref as entity_ref, db_entities.type, db_entities.title, db_identifiers.label, db_identifiers.value from db_entities, db_identifiers left join db_entities using(id, entity_ref);

select db_identifiers.entity_ref as entity_ref, db_entities.type, db_entities.title, db_identifiers.label, db_identifiers.value 
        from db_entities, db_identifiers left join db_entities using(id) where id = entity_ref;
        
select * from db_entities, db_identifiers, db_classifiers 
        where db_entities.id = db_identifiers.entity_ref 
          and db_entities.id = db_classifiers.entity_ref 
          and db_entities.id = 44000;        


select * from db_entities join db_references on db_entities.entity_id == db_references.entity_id where uuid = "18c65a86-74c5-4e5d-9203-b47dd368e986";
select * from db_entities join db_identifiers on db_entities.id == db_identifiers.entity_ref where uuid = "04e009ac-fe40-4ab8-8566-91e6304e6280";

select distinct * from db_entities as E, db_identifiers as I, db_references as R, db_classifiers as C 
        where E.uuid = "3d08f1ba-dcde-4d11-ae41-b7748309105f" 
        and (E.id == I.entity_ref and E.id == R.entity_ref and E.id == C.entity_ref);
        
select * from db_classifiers where entity_ref == 4572;
