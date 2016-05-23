from keras.models import Sequential
from keras.layers import Dense, Dropout, Activation, Lambda, Merge, Flatten
from keras.layers import Embedding
from keras.layers import Convolution1D, MaxPooling1D
from keras import backend as K
import numpy

import csv

INPUT_LENGTH = 29

V1_HIDDEN_SIZE = 4
EMBEDDING_SIZE = 256
NB_CONV_FILTER1 = 64
FILTER_LENGTH = 3
HIDDEN_DIMS = 64

BATCH_SIZE = 32
NB_EPOCH = 5
#LOSS_FUNCTION = 'binary_crossentropy'
LOSS_FUNCTION = 'mean_squared_error'

def load_data(filename, repair_only=False):
  values1 = []
  values2 = []
  labels = []
  with open(filename) as f:
    line = 0
    for row in csv.reader(f):
      line += 1
      i_row = [float(x) for x in row]
      if repair_only and i_row[-2] == 1:
        continue
      values1.append(i_row[1:5])
      values2.append(i_row[5:-1])
      labels.append(i_row[-1])
  values1 = numpy.asarray(values1)
  values2 = numpy.asarray(values2)
  values2 = values2.reshape(values2.shape + (1,))
  return values1, values2, labels

V1_train, V2_train, L_train = load_data('data/shuffled-train.csv')
V1_test, V2_test, L_test = load_data('data/shuffled-ai.csv', repair_only=True)

model1 = Sequential()
model1.add(Dense(V1_HIDDEN_SIZE, input_dim=4))

model2 = Sequential()

# we add a Convolution1D, which will learn nb_filter
# word group filters of size filter_length:
model2.add(Convolution1D(nb_filter=NB_CONV_FILTER1,
                         filter_length=FILTER_LENGTH,
                         input_dim=1,
                         input_length=INPUT_LENGTH,
                         border_mode='valid',
                         activation='relu'))

model2.add(MaxPooling1D(pool_length=2))
model2.add(Dropout(0.2))

model2.add(Flatten())
model2.add(Dense(HIDDEN_DIMS, activation='relu'))

model = Sequential()
model.add(Merge([model1, model2], mode='concat'))
# We add a vanilla hidden layer:
model.add(Dense(HIDDEN_DIMS))
model.add(Dropout(0.2))
model.add(Activation('relu'))

# We project onto a single unit output layer, and squash it with a sigmoid:
model.add(Dense(1))
model.add(Activation('sigmoid'))

model.compile(loss=LOSS_FUNCTION,
              optimizer='adam',
              metrics=['accuracy'])

model.fit([V1_train, V2_train], L_train,
          batch_size=BATCH_SIZE,
          nb_epoch=NB_EPOCH,
          validation_data=([V1_test, V2_test], L_test))

score, acc = model.evaluate([V1_test, V2_test], L_test, batch_size=BATCH_SIZE)
print('Test score:', score)
print('Test accuracy:', acc)
