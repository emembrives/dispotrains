days <- read.table('ttl.csv', header=FALSE, sep=',')
d_sncf <- days[which(days$V1=='Disponible' & days$V2 == 'sncf'), ]
d_ratp <- days[which(days$V1=='Disponible' & days$V2 == 'ratp'), ]
hs_ratp <- days[which(days$V1=='Hors service' & days$V2 == 'ratp'), ]
hs_sncf <- days[which(days$V1=='Hors service' & days$V2 == 'sncf'), ]
d1 <- density(d_ratp$V3/24.0, from=0, to=182)
d2 <- density(d_sncf$V3/24.0, from=0, to=182)

png("ttl.png", width=800, height=600)
plot(d2, col="red", main="Temps avant panne des ascenseurs STIF", xlab="Temps (jours)", ylab = "Probabilité")
lines(d1, col="blue")
abline(v=seq(0, 182, 10), lty=2)
legend("topright", title="Réseau", c("SNCF", "RATP"), fill=c("red", "blue"), bg="white")
dev.off()

d1 <- density(hs_ratp$V3, from=0, to=168, adjust=3)
d2 <- density(hs_sncf$V3, from=0, to=168, adjust=3)
png("ttd.png", width=800, height=600)
plot(d1, col="blue", main="Temps avant réparation des ascenseurs STIF", xlab="Temps (heures)", ylab = "Probabilité")
lines(d2, col="red")> abline(v=seq(0, 168, 24), lty=2)
legend("topright", title="Réseau", c("SNCF", "RATP"), fill=c("red", "blue"), bg="white")
dev.off()

