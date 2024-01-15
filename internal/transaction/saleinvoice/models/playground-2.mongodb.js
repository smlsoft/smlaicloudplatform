
const database = 'dedepos';
const collection = 'transactionSaleInvoice';

// The current database to use.
use(database);


db.transactionSaleInvoice.find({"details.manufacturerguid": {$exists:false}}).forEach(function(doc)  { 
    
    if (doc.details == null) {
        return;
    }

    for (let i = 0; i < doc.details.length; i++) {
        let pdt =  db.productBarcodes.findOne({shopid: doc.shopid, barcode: doc.details[i].barcode, "deletedat": {$exists: false}});
        
        if (pdt != null) {
            let tempManufacturerguid = pdt.manufacturerguid == null ? "" : pdt.manufacturerguid;
            doc.details[i].manufacturerguid = tempManufacturerguid;
        }else{
            doc.details[i].manufacturerguid = "";
        }
        
    }
    db.transactionSaleInvoice.updateOne({_id: doc._id}, {$set: {details: doc.details}});
}
);
