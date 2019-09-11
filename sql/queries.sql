select * from db_entities, db_identifiers where entity_id = "https://mart.getty.edu/museum/collection/person/da88f653-216e-4dd3-9816-eeeff36c81f3";

select * from db_entities join db_references on db_entities.entity_id == db_references.entity_id where uuid = "18c65a86-74c5-4e5d-9203-b47dd368e986";
select * from db_entities join db_identifiers on db_entities.id == db_identifiers.entity_ref where uuid = "04e009ac-fe40-4ab8-8566-91e6304e6280";

select distinct * from db_entities as E, db_identifiers as I, db_references as R, db_classifiers as C 
        where E.uuid = "3d08f1ba-dcde-4d11-ae41-b7748309105f" 
        and (E.id == I.entity_ref and E.id == R.entity_ref and E.id == C.entity_ref);
        
select distinct * from db_classifiers where entity_ref == 4572;
