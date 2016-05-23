days <- read.table('ttl.csv', header=FALSE, sep=',')
d_sncf <- days[which(days$V1=='Disponible' & days$V2 == 'sncf'), ]
d_ratp <- days[which(days$V1=='Disponible' & days$V2 == 'ratp'), ]
hs_ratp <- days[which(days$V1=='Hors service' & days$V2 == 'ratp'), ]
hs_sncf <- days[which(days$V1=='Hors service' & days$V2 == 'sncf'), ]

summary(d_sncf$V3)
summary(d_ratp$V3)
summary(hs_sncf$V3)
summary(hs_ratp$V3)

wilcox.test(d_sncf$V3, d_ratp$V3, conf.int = TRUE, conf.level = 0.99)
wilcox.test(hs_sncf$V3, hs_ratp$V3, conf.int = TRUE, conf.level = 0.99)
